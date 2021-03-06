package http

import (
	"battleship/config"
	"battleship/controllers"
	"battleship/di"
	"battleship/middlewares"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/ziflex/lecho/v2"
	"net/http"
)

func StartHttpServer() {
	e := echo.New()
	e.HideBanner = true
	setHttpMiddlewares(e)
	setHttpEndpoints(e)
	e.Debug = true
	e.Logger = lecho.From(log.Logger)
	httpConfig := &http.Server{
		Addr: fmt.Sprintf(":%s", config.C.HttpPort),
	}
	log.Fatal().Err(e.StartServer(httpConfig)).Msg("")
}

func setHttpMiddlewares(e *echo.Echo) {
	e.Use(middleware.BodyDump(middlewares.BodyDumper))
	e.Use(middlewares.LogMiddleware())
	if config.C.Cors.Domain == "*" {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			//MaxAge:       86400,
			//AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "cache-control"},
		}))
	} else {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{config.C.Cors.Domain},
			MaxAge:       86400,
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "cache-control"},
		}))
	}
}

var (
	gameController = di.CreateGameController()
	userController = di.CreateUserController()
	socketHandler  = di.CreateSocketHandler()
)

func setHttpEndpoints(e *echo.Echo) {
	e.GET("/api/v1/check-health", controllers.CheckHealth)
	e.GET("/api/v1/error", controllers.Error)
	e.POST("/api/v1/game", gameController.CreateGame)
	e.POST("/api/v1/game/join", gameController.JoinGame)
	e.POST("/api/v1/game/submit-ships", gameController.SubmitShipsLocations)
	e.POST("/api/v1/game/change-turn", gameController.ChangeTurn)
	e.POST("/api/v1/game/move-ship", gameController.MoveShip)
	e.POST("/api/v1/game/reveal", gameController.RevealEnemyFields)
	e.POST("/api/v1/game/explode", gameController.Explode)
	e.GET("/api/v1/game/:game_id", gameController.GetGame)
	e.POST("/api/v1/user", userController.CreateUser)
	e.GET("/api/v1/user/:user_id", userController.GetUser)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/socket", socketHandler.CreateSocket)
}
