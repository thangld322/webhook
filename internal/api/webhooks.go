package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"webhook/internal/repository"
	"webhook/pkg"
)

type WebhookInterface interface {
	Create(c *gin.Context)
}

type WebhookController struct {
	repo         repository.WebhookInterface
	cacheService *pkg.Cache
	logger       *logrus.Entry
}

func NewWebhookController(repo repository.WebhookInterface, cacheService *pkg.Cache) WebhookInterface {
	return &WebhookController{
		repo:         repo,
		cacheService: cacheService,
		logger: pkg.Logger.WithFields(logrus.Fields{
			"entity": "webhook",
		}),
	}
}

func (ctrl *WebhookController) Create(c *gin.Context) {
}
