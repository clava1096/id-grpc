package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Origin struct {
	gorm.Model
	Uuid uuid.UUID
	Link string `json:"link"`
}
