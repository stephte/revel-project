package models

import (
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type BaseModel struct {
	gorm.Model
	Key						uuid.UUID	`gorm:"type:uuid;uniqueIndex;not null;default:uuid_generate_v4()"`
}
