package repository

import "gorm.io/gorm"

type SubscriberInterface interface {
	Create(any) error
}

type Subscriber struct {
	db *gorm.DB
}

func NewSubscriber(database *gorm.DB) SubscriberInterface {
	return &Subscriber{db: database}
}

func (r *Subscriber) Create(_ any) error {

	return nil
}
