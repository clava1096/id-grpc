package models

import (
	"gorm.io/gorm"
	"id-backend-grpc/internal/models/consts"
	"time"
)

// TODO разбить на слой DTO. сервис знает о модели - плохо!

type User struct {
	gorm.Model
	Name             string           `json:"name"`
	Password         string           `json:"password"`
	Email            string           `gorm:"uniqueIndex:userEmail"`
	RegistrationDate time.Time        `json:"registrationDate"`
	UserRole         consts.UsersRole `gorm:"default:user"`
	ConfirmEmail     bool             `gorm:"default:false"`
	TelegramId       int64            `gorm:"default:0"`
	VkId             int64            `gorm:"default:0"`
	Salt             string           `gorm:"default:''"`
}

type VKIDUser struct {
	ID        string `json:"user_id"`
	Avatar    string `json:"avatar,omitempty"`
	Birthday  string `json:"birthday,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Sex       int    `json:"sex,omitempty"`
	Verified  bool   `json:"verified"`
}

type VKUser struct {
	User VKIDUser `json:"user"`
}

type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

type CreateUser struct {
	Email      string
	Password   string
	Name       string
	VKid       int64
	TelegramID int64
}

type UserInfo struct {
	Email            string           `json:"email"`
	Name             string           `json:"name"`
	RegistrationDate time.Time        `json:"registrationDate"`
	Avatar           string           `json:"avatar"`
	TelegramID       int64            `json:"telegram_id"`
	VkId             int64            `json:"vk_id"`
	UserRole         consts.UsersRole `json:"user_role"`
}

func NewUser(email, password, salt, name string) *User {
	user := &User{
		Email:            email,
		Password:         password,
		Name:             name,
		RegistrationDate: time.Now(),
		Salt:             salt,
	}
	return user
}
