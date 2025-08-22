package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"id-backend-grpc/internal/models"
	"net/http"
	"net/mail"
	"strconv"
)

type registerRequest struct {
	Email        string               `json:"email"`
	Password     string               `json:"password"`
	Name         string               `json:"name"`
	RedirectLink string               `json:"redirect_link"`
	TelegramUser *models.TelegramUser `json:"telegram_user"`
	VKUser       *models.VKIDUser     `json:"vk_user"`
}

func (is *IdentityService) Registration(c echo.Context) error {
	var r registerRequest
	if err := c.Bind(&r); err != nil && valid(&r) { // TODO подключи валидатор ПЕРЕПИШИ ЕГО
		return ErrorResponse(c, "Bad request", 400)
	}

	if r.VKUser != nil && r.TelegramUser != nil {
		return ErrorResponse(c, "Bad Request", 400)
	}

	if _, err := is.userService.GetUserByEmail(r.Email); err == nil { // может быть косяк
		return ErrorResponse(c, "Conflict", 409)
	}

	ctx := c.Request().Context()
	u, err := is.userService.CreateUser(ctx, &models.CreateUser{ // может получится так, что анонимные функции не успеют проверить тг и вк
		Email:    r.Email,
		Password: r.Password, // TODO не проверяешь что тебе пришло кроме пустой строки, используй валидатор ебло
		Name:     r.Name,
		TelegramID: func() int64 {
			if r.TelegramUser == nil {
				return 0
			}
			return r.TelegramUser.ID
		}(),
		VKid: func() int64 {
			if r.VKUser == nil {
				return 0
			}
			id, _ := strconv.ParseInt(r.VKUser.ID, 10, 64)
			return id
		}(),
	})
	if err != nil {
		return ErrorResponse(c, "Internal Server Error", 500)
	}
	s, _ := fmt.Printf("hello, %s", u.Name)
	return c.JSON(http.StatusOK, s)
}

func valid(r *registerRequest) bool {
	_, err := mail.ParseAddress(r.Email) // false when address invalid
	return !(err == nil) || r.Password == "" || r.Email == ""
}
