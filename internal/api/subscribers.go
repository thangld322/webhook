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
	repo   repository.SubscriberInterface
	logger *logrus.Entry
}

func NewSubscriberController(repo repository.SubscriberInterface) SubscriberInterface {
	return &SubscriberController{
		repo: repo,
		logger: pkg.Logger.WithFields(logrus.Fields{
			"entity": "subscriber",
		}),
	}
}

func (ctrl *SubscriberController) Create(c *gin.Context) {
}
