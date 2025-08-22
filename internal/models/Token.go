package models

import (
	"id-backend-grpc/internal/models/consts"
	"time"
)

type Token struct {
	Uuid      string `gorm:"unique"`
	Ipaddress string
	UserID    uint
	User      User             `gorm:"foreignKey:UserID"`
	Type      consts.TokenType `gorm:"default:2"` // token.authorization
	ExpiresAt time.Time        `json:"expires_in"`
}

type CreateToken struct {
	Ipaddress string
	UserID    uint
	Type      consts.TokenType
	ExpiresAt time.Duration
}

func TokenResponse(email string, token *Token) JsonTokenResponse {
	return JsonTokenResponse{
		Token: token.Uuid,
		Email: email,
		Type:  token.Type,
	}
}

type JsonTokenResponse struct {
	Email string           `json:"email"`
	Token string           `json:"token"`
	Type  consts.TokenType `json:"token_type"`
}
