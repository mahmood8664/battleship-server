package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type CheckHealthController struct {
}

func NewCheckHealthController() *CheckHealthController {
	return &CheckHealthController{}
}

func (r *CheckHealthController) CheckHealth() func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, "OK")
	}
}
