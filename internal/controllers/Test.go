package controllers

import "github.com/labstack/echo/v4"

func TestHandler(c echo.Context) error {
	return c.String(200, "OK.")
}
