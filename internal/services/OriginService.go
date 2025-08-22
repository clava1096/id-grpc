package services

import (
	"id-backend-grpc/internal/models"
	"id-backend-grpc/internal/repositories"
)

type OriginService struct {
	repo repositories.OriginRepository
}

func NewOriginService(repo repositories.OriginRepository) *OriginService {
	return &OriginService{
		repo: repo,
	}
}

func (s *OriginService) CreateOrigin(o *models.Origin) (*models.Origin, error) {
	return o, s.repo.Create(o)
}
