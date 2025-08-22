package controllers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"id-backend-grpc/internal/models"
	"id-backend-grpc/internal/models/consts"
	"net/http"
)

type originUUID struct {
	Link string `json:"link"`
}

func (is *IdentityService) ValidateOriginUUID(c echo.Context) error {
	t, _ := is.tokenService.GetToken(GetToken(c))
	if t.Type != consts.TokenAccess {
		return ErrorResponse(c, "Locked", http.StatusLocked)
	}
	header := c.Request().Header.Get("X-Origin-UUID")
	link, err := uuid.Parse(header)
	if err != nil {
		return ErrorResponse(c, "UUID is invalid", http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, models.JsonLinkResponse{Link: link.String()})
}

func (is *IdentityService) CreateOriginUUID(c echo.Context) error {
	t, _ := is.tokenService.GetToken(GetToken(c))
	if u, _ := is.userService.GetUserByID(t.UserID); u.UserRole != consts.Admin {
		return ErrorResponse(c, "Locked", http.StatusLocked)
	}
	var r originUUID
	if err := c.Bind(r); err != nil {
		return ErrorResponse(c, "Bad request", http.StatusBadRequest)
	}
	o, err := is.originService.CreateOrigin(&models.Origin{
		Uuid: uuid.New(),
		Link: r.Link,
	})
	if err != nil {
		return ErrorResponse(c, "Internal Server Error", http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, o)
}
