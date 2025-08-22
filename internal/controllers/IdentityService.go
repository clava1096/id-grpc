package controllers

import (
	"github.com/labstack/echo/v4"
	"id-backend-grpc/internal/services"
)

type IdentityService struct {
	tokenService  *services.TokenService
	userService   *services.UserService
	originService *services.OriginService
}

func NewIdentityService(tokenService *services.TokenService,
	userService *services.UserService,
	originService *services.OriginService) *IdentityService {
	return &IdentityService{
		tokenService:  tokenService,
		userService:   userService,
		originService: originService,
	}
}

// TODO ВНИМАНИЕ, КОД СГЕНЕРИЛА НЕЙРОГАВЕХА!!!!
func GetToken(c echo.Context) string {
	// Пробуем получить токен из куки
	cookie, err := c.Cookie("Access-Token")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// Пробуем получить из query параметров
	token := c.QueryParam("token")
	if token != "" {
		return token
	}

	token = c.QueryParam("tokenQuery")
	if token != "" {
		return token
	}

	// Пробуем получить из form данных
	token = c.FormValue("token")
	if token != "" {
		return token
	}

	// Пробуем получить из заголовка Authorization
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader != "" {
		return authHeader // На случай, если токен без 'Bearer'
	}

	return ""
}
