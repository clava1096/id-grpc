package controllers

import (
	"github.com/labstack/echo/v4"
	"id-backend-grpc/internal/models"
	"id-backend-grpc/internal/models/consts"
	"net/http"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type verificationCode struct {
	Code string `json:"code"`
}

func (is *IdentityService) Authorization(c echo.Context) error {
	var r loginRequest
	if err := c.Bind(r); err != nil {
		return ErrorResponse(c, "Bad request", 400)
	}
	u, err := is.userService.GetUserByEmail(r.Email)
	if err != nil {
		return ErrorResponse(c, "Not Found", 404)
	}
	if !u.ConfirmEmail {
		//TODO OTPService
	}
	if !is.userService.VerifyPasswordWithSalt(u.Password, r.Password, u.Salt) {
		return ErrorResponse(c, "Not Found", 404)
	}
	if !is.tokenService.AccessTokenIsExists(u.ID) {
		//TODO OTPService
		ctx := c.Request().Context()
		tokenTemporaryAccess, err := is.tokenService.GenerateToken(ctx, &models.CreateToken{
			Ipaddress: c.RealIP(),
			UserID:    u.ID,
			Type:      consts.TokenTemporaryAccess,
			ExpiresAt: consts.TokenTemporaryAccessLifeTime})
		if err != nil {
			return ErrorResponse(c, "Internal Server Error", 500)
		}
		return c.JSON(http.StatusOK, models.TokenResponse(u.Email, tokenTemporaryAccess))
	}
	token, _ := is.tokenService.GetTokenByUserIDAndType(u.ID, consts.TokenAccess)

	return c.JSON(http.StatusOK, models.TokenResponse(u.Email, token))
}

func (is *IdentityService) Logout(c echo.Context) error {
	token, err := is.tokenService.GetToken(GetToken(c))
	if err != nil {
		return ErrorResponse(c, "Access Denied", 403)
	}
	if token.Type != consts.TokenAccess {
		return ErrorResponse(c, "Locked", 423)
	}
	err = is.tokenService.Delete(token)
	if err != nil {
		return ErrorResponse(c, "Internal Server Error", 500)
	}
	return c.JSON(http.StatusOK, models.JsonResponse{Msg: "Logout successfully"})
}

func (is *IdentityService) Verify(c echo.Context) error {
	t, err := is.tokenService.GetToken(GetToken(c))
	if err != nil {
		return ErrorResponse(c, "Unauthorized", http.StatusUnauthorized)
	}
	if t.Type != consts.TokenTemporaryAccess {
		return ErrorResponse(c, "Invalid token type", http.StatusInternalServerError)
	}
	err = is.tokenService.Delete(t)
	if err != nil {
		return ErrorResponse(c, "Error while delete token temporary access", http.StatusInternalServerError)
	}
	var code verificationCode
	if err = c.Bind(code); err != nil {
		return ErrorResponse(c, "Bad request", 400)
	}
	if is.tokenService.VerifyOTPToken(code.Code) {
		ctx := c.Request().Context()
		u, _ := is.userService.GetUserByID(t.UserID)
		token, _ := is.tokenService.GenerateToken(ctx, &models.CreateToken{
			Ipaddress: c.RealIP(),
			UserID:    u.ID,
			Type:      consts.TokenAccess,
			ExpiresAt: consts.TokenAccessLifeTime})
		return c.JSON(http.StatusCreated, models.TokenResponse(u.Email, token))
	}
	return ErrorResponse(c, "Code is invalid or expired", http.StatusUnauthorized)
}
