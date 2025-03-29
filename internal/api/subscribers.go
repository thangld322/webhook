package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	"webhook/internal/model"
	"webhook/internal/repository"
	"webhook/pkg"
)

type SubscriberInterface interface {
	Create(c *gin.Context)
}

type SubscriberController struct {
	repo         repository.SubscriberInterface
	producer     *pkg.Producer
	cacheService *pkg.Cache
	logger       *logrus.Entry
}

func NewSubscriberController(repo repository.SubscriberInterface, producer *pkg.Producer, cacheService *pkg.Cache) SubscriberInterface {
	return &SubscriberController{
		repo:         repo,
		producer:     producer,
		cacheService: cacheService,
		logger: pkg.Logger.WithFields(logrus.Fields{
			"entity": "subscriber",
		}),
	}
}

func (ctrl *SubscriberController) Create(c *gin.Context) {
	var subscriber *model.Subscriber
	var err error
	if err = c.Bind(&subscriber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if err = subscriber.GenerateID(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	subscriber.CreatedAt = time.Now()
	subscriber.UpdatedAt = subscriber.CreatedAt

	err = ctrl.repo.Create(subscriber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	// Create Event
	event := model.WebhookEvent{
		TenantID:   subscriber.TenantID,
		EventName:  "subscriber.created",
		EventTime:  time.Now(),
		Subscriber: subscriber,
	}
	eventByte, err := json.Marshal(event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	// Choose topic
	var topic string
	count, err := ctrl.cacheService.Get(fmt.Sprintf(model.TenantQueueCount, subscriber.TenantID))
	if err != nil {
		if err.Error() == "redis: nil" {
			_, err = ctrl.cacheService.Set(fmt.Sprintf(model.TenantQueueCount, subscriber.TenantID), 0, 0)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
				return
			}
			count = 1
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}
	}
	if count > model.TenantQueueThreshold {
		topic = fmt.Sprintf(model.TenantTopic, subscriber.TenantID)
	} else {
		topic = fmt.Sprintf(model.TenantTopic, "common")
	}

	// Produce a message to the topic
	err = ctrl.producer.Produce(topic, eventByte)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	// Increase tenant queue count
	count, err = ctrl.cacheService.Incr(fmt.Sprintf(model.TenantQueueCount, subscriber.TenantID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	ctrl.logger.Printf("Topic: %s, tenant: %s, count: %d\n", topic, subscriber.TenantID, count)

	c.JSON(http.StatusOK, subscriber)
}
