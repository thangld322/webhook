package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Model struct {
	ID        string `gorm:"primary_key;type:varchar(100)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (model *Model) GenerateID() error {
	uv4, err := uuid.NewV4()
	if err != nil {
		return err
	}
	model.ID = uv4.String()
	return nil
}
