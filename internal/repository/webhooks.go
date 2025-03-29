package repository

import (
	"gorm.io/gorm"

	"webhook/internal/model"
)

type WebhookInterface interface {
	Create(webhook *model.Webhook) error
	GetByEvent(tenantID, event string) ([]model.Webhook, error)
	UpdateStatus(id string, status bool) error
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

func (r *Webhook) GetByEvent(tenantID, event string) ([]model.Webhook, error) {
	var webhooks []model.Webhook
	if err := r.db.Where("is_active = true AND tenant_id = ? AND ? = ANY(events)", tenantID, event).Find(&webhooks).Error; err != nil {
		return nil, err
	}

	return webhooks, nil
}

func (r *Webhook) UpdateStatus(id string, status bool) error {
	return r.db.Model(&model.Webhook{}).Where("id = ", id).Update("is_active", status).Error
}
