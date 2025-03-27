package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"webhook/internal/repository"
	"webhook/pkg"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, cacheService *pkg.Cache) {
	router.RedirectTrailingSlash = true

	// Init repo
	webhookRepo := repository.NewWebhook(db)
	subscriberRepo := repository.NewSubscriber(db)

	// Init controller
	webhookController := NewWebhookController(webhookRepo, cacheService)
	subscriberController := NewSubscriberController(subscriberRepo, cacheService)

	v1 := router.Group("/v1")
	{
		v1.POST("/webhooks", webhookController.Create)

		v1.POST("/subscribers", subscriberController.Create)
	}

}
