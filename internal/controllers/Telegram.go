package controllers

import (
	"github.com/labstack/echo/v4"
	"id-backend-grpc/internal/app/config"
	"id-backend-grpc/internal/models"
	"id-backend-grpc/internal/models/consts"
	"net/http"
	"time"
)

func (is *IdentityService) TelegramAuthCallback(c echo.Context) error {
	var r models.TelegramUser // TODO это дто слой долбаеб, убрать из файла моделей
	cfg, _ := config.LoadTelegramConfig()
	if err := c.Bind(r); err != nil || is.userService.ValidateTelegramData(&r, cfg.Token) {
		return ErrorResponse(c, "Bad Request", http.StatusBadRequest)
	}
	u, _ := is.userService.GetUserByID(uint(r.ID))
	if !is.tokenService.AccessTokenIsExists(u.ID) {
		// TODO OTPService
		ctx := c.Request().Context()
		token, _ := is.tokenService.GenerateToken(ctx, &models.CreateToken{
			Ipaddress: c.RealIP(),
			UserID:    u.ID,
			Type:      consts.TokenAccess,
			ExpiresAt: consts.TokenAccessLifeTime})
		return c.JSON(http.StatusCreated, models.TokenResponse(u.Email, token))
	}
	if time.Now().Unix()-r.AuthDate > 86400 {
		return ErrorResponse(c, "Unauthorized", http.StatusUnauthorized)
	}
	if !u.ConfirmEmail {
		//TODO OTPService
		return ErrorResponse(c, "You email is not confirm!", http.StatusUnauthorized)
	}
	token, _ := is.tokenService.GetTokenByUserIDAndType(u.ID, consts.TokenAccess)
	return c.JSON(http.StatusOK, models.TokenResponse(u.Email, token))
}

func (is *IdentityService) TelegramIntegration(c echo.Context) error {
	t, _ := is.tokenService.GetToken(GetToken(c))
	if t.Type != consts.TokenAccess {
		return ErrorResponse(c, "Locked", 423)
	}
	var r models.TelegramUser
	cfg, _ := config.LoadTelegramConfig()
	if err := c.Bind(r); err != nil || is.userService.ValidateTelegramData(&r, cfg.Token) {
		return ErrorResponse(c, "Bad request", 400)
	}
	u, _ := is.userService.GetUserByID(t.UserID)
	u = is.userService.SetTelegramID(u, r.ID)
	err := is.userService.SaveUser(u)
	if err != nil {
		return ErrorResponse(c, "Internal Server Error", http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, models.JsonResponse{Msg: "OK. Telegram account is integrated"})
}
