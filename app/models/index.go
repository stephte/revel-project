package models

import (
	"gorm.io/gorm"
	"github.com/google/uuid"
	"errors"
)

type BaseModel struct {
	gorm.Model
	Key						uuid.UUID	`gorm:"type:uuid;uniqueIndex;not null;default:uuid_generate_v4()"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.Key == nil {
		m.Key = uuid.New()
	}

	if !m.IsValid() {
		return errors.New("cannot create user")
	}

	return nil
}

func (m *BaseModel) BeforeSave(tx *gorm.DB) (err error) {
	if !m.IsValid() {
		return errors.New("Model invalid")
	}
}
