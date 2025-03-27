package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"webhook/pkg"

	"webhook/internal/repository"
)

type WebhookInterface interface {
	Create(c *gin.Context)
}

type WebhookController struct {
	repo   repository.WebhookInterface
	logger *logrus.Entry
}

func NewWebhookController(repo repository.WebhookInterface) WebhookInterface {
	return &WebhookController{
		repo: repo,
		logger: pkg.Logger.WithFields(logrus.Fields{
			"entity": "webhook",
		}),
	}
}

func (ctrl *WebhookController) Create(c *gin.Context) {
}
