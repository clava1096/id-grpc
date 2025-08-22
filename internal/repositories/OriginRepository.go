package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"id-backend-grpc/internal/app/config/connections"
	"id-backend-grpc/internal/models"
)

type OriginRepository interface {
	Create(o *models.Origin) error
}

type originRepository struct {
	db *gorm.DB
}

func NewOriginRepository() (OriginRepository, error) {
	db, err := connections.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create DB connection: %w", err)
	}
	return &originRepository{db: db}, nil
}

func (r originRepository) Create(o *models.Origin) error {
	return r.db.Create(&o).Error
}
