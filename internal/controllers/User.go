package controllers

import (
	"github.com/labstack/echo/v4"
	"id-backend-grpc/internal/models"
	"id-backend-grpc/internal/models/consts"
	"net/http"
	"os"
)

func (is *IdentityService) User(c echo.Context) error {
	token, _ := is.tokenService.GetToken(GetToken(c))
	if token.Type != consts.TokenAccess {
		return ErrorResponse(c, "Locked", http.StatusLocked)
	}
	user, _ := is.userService.GetUserByID(token.UserID)
	data := models.UserInfo{
		Email:            user.Email,
		Name:             user.Name,
		RegistrationDate: user.RegistrationDate,
		//Avatar:           config.GetConfig().Application.Host + "/api/1.0/user.get-avatar", // TODO неверно
		TelegramID: user.TelegramId,
		VkId:       user.VkId,
		UserRole:   user.UserRole,
	}
	return c.JSON(http.StatusOK, data)
}

func (is *IdentityService) UploadAvatar(c echo.Context) error {
	token, _ := is.tokenService.GetToken(GetToken(c))
	if token.Type != consts.TokenAccess {
		return ErrorResponse(c, "Locked", http.StatusLocked)
	}
	form, _ := c.MultipartForm()
	u, _ := is.userService.GetUserByID(token.UserID)
	if !is.userService.UploadAvatar(u, form) {
		return c.JSON(http.StatusInternalServerError, models.JsonResponse{Msg: "Error while upload avatar"})
	}
	return c.JSON(http.StatusCreated, models.JsonResponse{Msg: "OK. Avatar Uploaded."})
}

func (is *IdentityService) LoadAvatar(c echo.Context) error {
	token, _ := is.tokenService.GetToken(GetToken(c))
	if token.Type != consts.TokenAccess {
		return ErrorResponse(c, "Locked", http.StatusLocked)
	}
	u, _ := is.userService.GetUserByID(token.UserID)
	get := is.userService.GetAvatar(u)
	// TODO ПЕРЕПИСАТЬ НИЖЕ, НАВРЯД ЛИ РАБОТАЕТ
	defer os.Remove(get.Name())

	if get == nil {
		return ErrorResponse(c, "Internal Server Error", http.StatusInternalServerError)
	}
	return c.File(get.Name())
}
