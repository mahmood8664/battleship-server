package controllers

import (
	"battleship/battle_error"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CheckHealth() func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, "OK")
	}
}

func Error() func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		return battle_error.BadRequest1("This is Bad Request")
	}
}
