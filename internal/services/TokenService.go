package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"id-backend-grpc/internal/models"
	"id-backend-grpc/internal/models/consts"
	"id-backend-grpc/internal/repositories"
	"math/rand"
	"time"
)

type TokenService struct {
	repo repositories.TokenRepository
}

func NewTokenService(repo repositories.TokenRepository) *TokenService {
	return &TokenService{
		repo: repo,
	}
}

func (s *TokenService) GenerateToken(ctx context.Context, ct *models.CreateToken) (*models.Token, error) {
	token := &models.Token{
		Ipaddress: ct.Ipaddress,
		Uuid:      uuid.New().String(),
		ExpiresAt: time.Now().Add(ct.ExpiresAt),
		UserID:    ct.UserID,
		Type:      ct.Type,
	}
	if ct.Type == consts.TokenVerification {
		token.Uuid = generateOTPToken() // !!!!!!!! 	TODO отдельный сервис для отправки отп кодов, передавать этот токен туда по grpc
	}
	err := s.repo.Create(ctx, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *TokenService) GetTokenByUserIDAndType(UserID uint, Type consts.TokenType) (*models.Token, error) {
	return s.repo.GetTokenByUserIDAndType(UserID, Type)
}

func (s *TokenService) AccessTokenIsExists(userID uint) bool {
	return s.repo.IsExists(userID)
}

func (s *TokenService) GetToken(uuid string) (*models.Token, error) {
	return s.repo.GetToken(uuid)
}

func (s *TokenService) Delete(token *models.Token) error {
	return s.repo.Delete(token)
}

func (s *TokenService) VerifyOTPToken(code string) bool {
	t, _ := s.repo.GetTokenByTokenAndType(code, consts.TokenVerification)
	if time.Now().After(t.ExpiresAt) {
		return false
	}
	_ = s.repo.Delete(t)
	return true
}

func generateOTPToken() string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	return fmt.Sprintf("%06d", rng.Intn(1000000))
}
