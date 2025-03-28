package repository

import (
	"gorm.io/gorm"

	"webhook/internal/model"
)

type WebhookInterface interface {
	Create(webhook *model.Webhook) error
}

type Webhook struct {
	db *gorm.DB
}

func NewWebhook(database *gorm.DB) WebhookInterface {
	return &Webhook{db: database}
}

func (r *Webhook) Create(webhook *model.Webhook) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(webhook).Error
}
