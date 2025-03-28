package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

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
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	if err = webhook.GenerateID(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	webhook.CreatedAt = time.Now()
	webhook.UpdatedAt = webhook.CreatedAt

	err = ctrl.repo.Create(webhook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, webhook)
}
