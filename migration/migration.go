package migration

import (
	"gorm.io/gorm"

	"webhook/internal/model"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Webhook{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&model.Subscriber{}); err != nil {
		return err
	}

	return nil
}
