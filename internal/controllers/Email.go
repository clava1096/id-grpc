package controllers

import (
	"github.com/labstack/echo/v4"
	"id-backend-grpc/internal/models"
	"id-backend-grpc/internal/models/consts"
	"log"
	"net/http"
	"time"
)

func (is *IdentityService) ConfirmEmail(c echo.Context) error {
	//t, err := is.tokenService.GetToken(GetToken(c))
	t, err := is.tokenService.GetToken(c.QueryParam("tokenQuery"))
	if err != nil || t.Type != consts.TokenConfirmEmail {
		return ErrorResponse(c, "This token was activated or have invalid type!", http.StatusForbidden)
	}
	if t.ExpiresAt.Before(time.Now()) {
		//TODO OTPService
		return c.JSON(http.StatusCreated, models.JsonResponse{Msg: "Please, confirm your email. А new letter has been sent on your email"})
	}
	err = is.userService.ConfirmEmailUserByID(t.UserID)
	if err != nil {
		log.Print(err.Error())
		return ErrorResponse(c, "Internal Server Error", http.StatusInternalServerError)
	}
	_ = is.tokenService.Delete(t)
	return c.JSON(http.StatusOK, models.JsonResponse{Msg: "Ok. Email confirmed"})
}
