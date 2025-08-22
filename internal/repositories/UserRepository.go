package repositories

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"id-backend-grpc/internal/app/config/connections"
	"id-backend-grpc/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	Delete(u *models.User) error
	Save(u *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByTelegramID(TelegramID uint) (*models.User, error) // TODO используется беззнаковый тип что верно, но в модели обычный
	GetUserByVKID(VKID uint) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() (UserRepository, error) {
	db, err := connections.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create DB connection: %w", err)
	}
	return &userRepository{db: db}, nil
}

func (r *userRepository) Create(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepository) Delete(u *models.User) error {
	return r.db.Delete(&u).Error
}

func (r *userRepository) Save(u *models.User) error {
	return r.db.Save(&u).Error
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) GetUserByTelegramID(TelegramID uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("telegram_id = ?", TelegramID).First(&user).Error
	return &user, err
}

func (r *userRepository) GetUserByVKID(VKID uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("vk_id = ?", VKID).First(&user).Error
	return &user, err
}

func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}
