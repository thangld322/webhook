package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"sync"

	"webhook/internal/repository"
	"webhook/internal/webhook"
	"webhook/pkg"
)

var notifierCmd = &cobra.Command{
	Use:   "notifier",
	Short: "Run the Webhook Notifier",
	RunE:  runNotifierCmd,
}

func runNotifierCmd(cmd *cobra.Command, args []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	// Init PostgreSQL
	var err error
	var orm *gorm.DB
	orm, err = NewDBConnection()
	if err != nil {
		pkg.Logger.WithError(err).Error("failed to init db connection")
		return err
	}
	webhookRepo := repository.NewWebhook(orm)

	// Init redis
	cacheService, err := pkg.NewCacheCache(ctx, redisAddress, redisPassword, redisPoolSize, redisMinIdleConns, redisDB)
	if err != nil {
		pkg.Logger.WithError(err).Error("failed to init redis client")
		return err
	}

	// Init kafka
	kafkaConsumer, err := pkg.NewConsumer(kafkaBroker)
	if err != nil {
		pkg.Logger.WithError(err).Error("failed to init kafka client")
		return err
	}

	s := &sync.WaitGroup{}

	notifier := webhook.NewNotifier(ctx, webhookRepo, kafkaConsumer, cacheService, 1, s)
	if err := notifier.Start(); err != nil {
		return err
	}
	<-ctx.Done()
	s.Wait()

	return nil
}
