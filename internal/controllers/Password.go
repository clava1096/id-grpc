package controllers

import (
	"github.com/labstack/echo/v4"
	"id-backend-grpc/internal/models"
	"id-backend-grpc/internal/models/consts"
	"net/http"
	"time"
)

type resetRequest struct {
	Email string `json:"email"`
	Link  string `json:"link"`
}

type resetPassword struct { // TODO не нравится название, переделать
	Password string `json:"password"`
}
type changePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (is *IdentityService) ResetPassword(c echo.Context) error {
	var r resetRequest
	if err := c.Bind(r); err != nil {
		return ErrorResponse(c, "Bad request", http.StatusBadRequest)
	}
	u, err := is.userService.GetUserByEmail(r.Email)
	if err != nil {
		return ErrorResponse(c, "Not Found", http.StatusNotFound)
	}
	ctx := c.Request().Context()
	_, _ = is.tokenService.GenerateToken(ctx, &models.CreateToken{
		Ipaddress: c.RealIP(),
		UserID:    u.ID,
		Type:      consts.TokenResetPassword,
		ExpiresAt: consts.TokenResetPasswordLifeTime})
	//TODO OTPService
	return c.JSON(http.StatusOK, models.JsonResponse{Msg: "Password reset link has been sent to your email"})
}

func (is *IdentityService) SaveResetPassword(c echo.Context) error {
	var r resetPassword
	t, err := is.tokenService.GetToken(c.QueryParam("token"))
	if err != nil {
		return ErrorResponse(c, "Invalid token", http.StatusBadRequest)
	}
	if err = c.Bind(r); err != nil {
		return ErrorResponse(c, "Bad request", http.StatusBadRequest)
	}
	if t.ExpiresAt.Before(time.Now()) {
		return ErrorResponse(c, "You must a send new request for reset password", http.StatusRequestTimeout)
	}
	u, _ := is.userService.GetUserByID(t.UserID)
	u, _ = is.userService.EditPassword(u, r.Password)
	err = is.userService.SaveUser(u)
	if err != nil {
		return ErrorResponse(c, "Internal Server Error", http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, models.JsonResponse{Msg: "Password has been changed to a new one"})
}

func (is *IdentityService) ChangePassword(c echo.Context) error {
	t, _ := is.tokenService.GetToken(GetToken(c))
	if t.Type != consts.TokenAccess {
		return ErrorResponse(c, "Locked", http.StatusLocked)
	}
	var r changePassword
	if err := c.Bind(r); err != nil {
		return ErrorResponse(c, "Bad request", http.StatusBadRequest)
	}
	u, _ := is.userService.GetUserByID(t.UserID)
	if !is.userService.VerifyPasswordWithSalt(u.Password, u.Salt, r.OldPassword) {
		return ErrorResponse(c, "Bad request! Your old password doesn't match", http.StatusBadRequest)
	}
	u, _ = is.userService.EditPassword(u, r.NewPassword)
	err := is.userService.SaveUser(u)
	if err != nil {
		return ErrorResponse(c, "Internal Server Error", http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, models.JsonResponse{Msg: "Password has been changed to a new one"})
}
