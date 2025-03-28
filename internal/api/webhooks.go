package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"webhook/internal/model"

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
	var webhook *model.Webhook
	var err error
	if err = c.Bind(&webhook); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
}
