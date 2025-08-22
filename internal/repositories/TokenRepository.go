package repositories

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"id-backend-grpc/internal/app/config/connections"
	"id-backend-grpc/internal/models"
	"id-backend-grpc/internal/models/consts"
)

type TokenRepository interface {
	Create(ctx context.Context, t *models.Token) error
	Delete(t *models.Token) error
	Save(t *models.Token) error
	IsExists(UserID uint) bool
	GetTokenByUserIDAndType(UserID uint, Type consts.TokenType) (*models.Token, error)
	GetAllTokensByType(typeToken consts.TokenType) (tokens *[]models.Token, err error) // возвожно добавить сюда контекст
	GetToken(uuid string) (*models.Token, error)                                       // зачем? если можно проверить только его существование
	GetTokenByTokenAndType(token string, typeToken consts.TokenType) (*models.Token, error)
	Close() error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository() (TokenRepository, error) {
	db, err := connections.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create DB connection: %w", err)
	}
	return &tokenRepository{db: db}, nil
}

func (r *tokenRepository) Create(ctx context.Context, t *models.Token) error {
	return r.db.WithContext(ctx).Create(&t).Error
}

func (r *tokenRepository) Delete(t *models.Token) error {
	return r.db.Delete(&t).Error
}

func (r *tokenRepository) Save(t *models.Token) error {
	return r.db.Save(&t).Error
}

func (r *tokenRepository) IsExists(UserID uint) bool {
	var isExists bool
	r.db.Model(&models.Token{}).
		Select("count(*) > 0").
		Where("user_id = ? and type = ?", UserID, consts.TokenAccess).
		Find(&isExists)
	return isExists
}

func (r *tokenRepository) GetTokenByUserIDAndType(UserId uint, Type consts.TokenType) (*models.Token, error) {
	var token models.Token
	err := r.db.Where("user_id = ? and type = ?", UserId, Type).First(&token).Error
	return &token, err
}

func (r *tokenRepository) GetAllTokensByType(typeToken consts.TokenType) (tokens *[]models.Token, err error) {
	err = r.db.Model(&models.Token{}).Where("type = ?", typeToken).Find(&tokens).Error
	return tokens, err
}

func (r *tokenRepository) GetToken(uuid string) (*models.Token, error) {
	var ExToken models.Token
	err := r.db.Where("token = ?", uuid).First(&ExToken).Error
	return &ExToken, err
}

func (r *tokenRepository) GetTokenByTokenAndType(token string, typeToken consts.TokenType) (*models.Token, error) {
	var exToken models.Token
	err := r.db.Where("token = ? and type = ?", token, typeToken).First(&exToken).Error
	return &exToken, err
}

func (r *tokenRepository) Close() error {
	return connections.Close(r.db)
}
