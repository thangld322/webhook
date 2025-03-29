package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"webhook/internal/api"
	"webhook/pkg"
)

func Start(ctx context.Context) error {
	logger := pkg.Logger.WithFields(logrus.Fields{
		"type": "server",
	})

	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	// Init PostgreSQL
	var err error
	var orm *gorm.DB
	orm, err = NewDBConnection()
	if err != nil {
		logger.WithError(err).Error("failed to init db connection")
		return err
	}

	// Init redis
	cacheService, err := pkg.NewCacheCache(ctx, redisAddress, redisPassword, redisPoolSize, redisMinIdleConns, redisDB)
	if err != nil {
		logger.WithError(err).Error("failed to init redis client")
		return err
	}

	// Init kafka
	kafkaProducer, err := pkg.NewProducer(kafkaBroker)
	if err != nil {
		logger.WithError(err).Error("failed to init kafka client")
		return err
	}

	// Gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	// Register HTTP route handlers.
	api.RegisterRoutes(router, orm, kafkaProducer, cacheService)

	// Create new HTTP server instance.
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "localhost", 8080),
		Handler: router,
	}

	logger.Debugf("http: successfully initialized")

	// Context with cancel
	cctx, cancelFunc := context.WithCancel(ctx)

	// Start HTTP server.
	go func(ctx context.Context) {
		go func() {
			logger.Infof("http: starting web server at %s", server.Addr)

			if err = server.ListenAndServe(); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					logger.Info("http: web server shutdown complete")
				} else {
					logger.Errorf("http: web server closed unexpect: %s", err)
				}
			}
		}()

		// Graceful HTTP server shutdown.
		<-ctx.Done()
		logger.Info("http: shutting down web server")
		err = server.Close()
		if err != nil {
			logger.Errorf("http: web server shutdown failed: %v", err)
		}
	}(cctx)

	// Handle sigterm and await termChan signal
	terminateChan := make(chan os.Signal)
	signal.Notify(terminateChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) //nolint

	<-terminateChan // Blocks here until interrupted

	// Handle shutdown
	logger.Println("*********************************\nShutdown signal received\n*********************************")
	cancelFunc() // Signal cancellation to context.Context

	time.Sleep(2 * time.Second)

	return nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
