package controllers

import (
	"battleship/dto"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CheckHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, "OK")
}

func Error(ctx echo.Context) error {
	return dto.BadRequest1("This is Bad Request")
}
