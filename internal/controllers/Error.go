package controllers

import (
	"github.com/labstack/echo/v4"
	"id-backend-grpc/internal/models"
)

func ErrorResponse(c echo.Context, err string, code int) error {
	return c.JSON(code, models.ErrorResponse{Error: err})
}
