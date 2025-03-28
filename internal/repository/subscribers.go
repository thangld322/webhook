package repository

import (
	"gorm.io/gorm"

	"webhook/internal/model"
)

type SubscriberInterface interface {
	Create(subscriber *model.Subscriber) error
}

type Subscriber struct {
	db *gorm.DB
}

func NewSubscriber(database *gorm.DB) SubscriberInterface {
	return &Subscriber{db: database}
}

func (r *Subscriber) Create(subscriber *model.Subscriber) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(subscriber).Error
}
