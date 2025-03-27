package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"webhook/pkg"

	"webhook/internal/repository"
)

type SubscriberInterface interface {
	Create(c *gin.Context)
}

type SubscriberController struct {
	repo         repository.SubscriberInterface
	cacheService *pkg.Cache
	logger       *logrus.Entry
}

func NewSubscriberController(repo repository.SubscriberInterface, cacheService *pkg.Cache) SubscriberInterface {
	return &SubscriberController{
		repo:         repo,
		cacheService: cacheService,
		logger: pkg.Logger.WithFields(logrus.Fields{
			"entity": "subscriber",
		}),
	}
}

func (ctrl *SubscriberController) Create(c *gin.Context) {
}
