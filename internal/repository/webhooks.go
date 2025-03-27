package repository

import "gorm.io/gorm"

type WebhookInterface interface {
	Create(any) error
}

type Webhook struct {
	db *gorm.DB
}

func NewWebhook(database *gorm.DB) WebhookInterface {
	return &Webhook{db: database}
}

func (r *Webhook) Create(_ any) error {

	return nil
}
