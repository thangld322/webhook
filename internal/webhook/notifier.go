package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"sync"
	"time"

	"webhook/internal/model"
	"webhook/internal/repository"
	"webhook/pkg"
)

type Notifier struct {
	ctx          context.Context
	repo         repository.WebhookInterface
	consumer     *pkg.Consumer
	cacheService *pkg.Cache

	numConsumer int
	jobsStream  chan *kafka.Message
	wg          *sync.WaitGroup
	logger      *logrus.Entry
}

func NewNotifier(ctx context.Context, repo repository.WebhookInterface, consumer *pkg.Consumer,
	cacheService *pkg.Cache, numConsumer int, wg *sync.WaitGroup) *Notifier {
	return &Notifier{
		ctx:          ctx,
		repo:         repo,
		consumer:     consumer,
		cacheService: cacheService,
		numConsumer:  numConsumer,
		jobsStream:   make(chan *kafka.Message, numConsumer),
		wg:           wg,
		logger: pkg.Logger.WithFields(logrus.Fields{
			"service": "notifier",
		}),
	}
}

func (n *Notifier) allocateMessage() {
	n.logger.Info("start allocate message")
	// Subscribe to topics matching the regex pattern
	pattern := "^tenant-.*"
	err := n.consumer.Client.SubscribeTopics([]string{pattern}, nil)
	if err != nil {
		n.logger.WithError(err).WithField("topic", pattern).Error("subscribe topic error")
		return
	}
	go func() {
		for {
			select {
			case <-n.ctx.Done():
				n.logger.Info("stop allocate message")
				n.Stop()
				return
			default:
				ev := n.consumer.Client.Poll(100)
				switch e := ev.(type) {
				case *kafka.Message:
					n.logger.Printf("Received message on %s: %s\n", e.TopicPartition, string(e.Value))
					n.jobsStream <- ev.(*kafka.Message)
				case kafka.Error:
					n.logger.Printf("Error: %v\n", e)
				}
			}
		}
	}()
}

func (n *Notifier) Start() error {
	host, err := os.Hostname()
	if err != nil {
		return err
	}
	n.allocateMessage()
	for i := 1; i <= n.numConsumer; i++ {
		name := fmt.Sprintf("webhook worker: %s-%d", host, i)
		consumer := newWebhookWorker(n, name)
		n.wg.Add(1)
		go consumer.Handle()
	}
	return nil
}

func (n *Notifier) Stop() error {
	if err := n.consumer.Client.Close(); err != nil {
		return err
	}
	close(n.jobsStream)
	return nil
}

type webhookWorker struct {
	*Notifier
	name string
}

func newWebhookWorker(n *Notifier, name string) *webhookWorker {
	return &webhookWorker{
		n,
		name,
	}
}

func (c *webhookWorker) Handle() {
	logger := c.logger.WithField("function", "Handle").WithField("consumer_name", c.name)
	defer c.wg.Done()
	for {
		select {
		case job, ok := <-c.jobsStream:
			if !ok {
				logger.Info("job stream already is closed")
				return
			}
			if err := c.execute(job); err != nil {
				logger.WithError(err).Error("handle job error")
			} else {
				logger.Info("handle job success")
			}
		case <-c.ctx.Done():
			logger.Info("closed consumer")
			return
		}
	}
}

func (c *webhookWorker) execute(job *kafka.Message) error {
	logger := c.logger.WithField("function", "execute").WithField("consumer_name", c.name)
	defer func() {
		_, err := c.consumer.Client.CommitMessage(job)
		if err != nil {
			logger.WithError(err).Error("commit message error")
		}
	}()

	var event model.WebhookEvent
	if err := json.Unmarshal(job.Value, &event); err != nil {
		return err
	}

	webhooks, err := c.repo.GetByEvent(event.TenantID, event.EventName)
	if err != nil {
		logger.WithError(err).Error("get webhooks by event failed")
		return err
	}

	s := &sync.WaitGroup{}
	for _, webhook := range webhooks {
		s.Add(1)
		go c.sendWebhook(job, webhook, s)
	}
	s.Wait()

	return nil
}

func (c *webhookWorker) sendWebhook(job *kafka.Message, webhook model.Webhook, s *sync.WaitGroup) {
	defer s.Done()
	// Create a new POST request with the payload
	req, err := http.NewRequest("POST", webhook.PostUrl, bytes.NewBuffer(job.Value))
	if err != nil {
		c.logger.WithError(err).WithField("webhook_id", webhook.ID).Error("create new request failed")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Retry
	backoffTime := time.Second
	retries := 0
	maxRetries := 5
	givingUp := false

	// Create an HTTP client and send the request
	client := &http.Client{}
	for {
		resp, err := client.Do(req)
		if err != nil || (resp != nil && resp.StatusCode > 299) {
			c.logger.WithError(err).WithField("webhook_id", webhook.ID).Error("send request failed")
			retries++
			if retries >= maxRetries {
				c.logger.Info("Max retries reached. Giving up on webhook:", webhook.ID)
				givingUp = true
				break
			}
			time.Sleep(backoffTime)
			backoffTime *= 2
			continue
		}
		break
	}

	if givingUp {
		err = c.repo.UpdateStatus(webhook.ID, false)
		if err != nil {
			c.logger.WithError(err).Error("update webhook status failed")
		}
	}
}
