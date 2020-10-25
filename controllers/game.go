package controllers

import (
	"battleship/battle_error"
	"battleship/dto"
	"battleship/service"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

type GameController interface {
	CreateGame() func(ctx echo.Context) error
	GetGame() func(ctx echo.Context) error
	JoinGame() func(ctx echo.Context) error
}

type GameControllerImpl struct {
	gameService service.GameService
}

func NewGameControllerImpl(gameService service.GameService) *GameControllerImpl {
	return &GameControllerImpl{
		gameService: gameService,
	}
}

// Create game
// @Summary Create game
// @Description create a new battleship game instance
// @Tags Game
// @Accept json
// @Produce json
// @Param request body dto.CreateGameRequest true "Create Game Request"
// @Success 200 {object} dto.CreateGameResponse "Game successfully created"
// @Router /api/v1/game [post]
func (r *GameControllerImpl) CreateGame() func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		request := new(dto.CreateGameRequest)
		if err := ctx.Bind(request); err != nil {
			log.Warn().Err(err).Msg("Bad request")
			return battle_error.BadRequest1(err.Error())
		}

		id, err := r.gameService.CreateGame(*request)
		if err != nil {
			log.Info().Str("userId", request.UserId).Err(err).Msg("cannot create game")
			return err
		}

		return ctx.JSON(http.StatusOK, dto.CreateGameResponse{
			GameId: id,
		})
	}
}

// Get game info
// @Summary Get game
// @Description Get game info
// @Tags Game
// @Accept json
// @Produce json
// @Param game_id path string true " "
// @Success 200 {object} dto.GameDto "Game successfully created"
// @Router /api/v1/game/{game_id} [get]
func (r *GameControllerImpl) GetGame() func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		gameId := ctx.Param("game_id")
		if gameId == "" {
			log.Warn().Msg("Bad request")
			return battle_error.BadRequest1("game_id must has value")
		}

		game, err := r.gameService.GetGame(gameId)
		if err != nil {
			log.Info().Str("gameId", gameId).Err(err).Msg("cannot get Game")
			return err
		}
		return ctx.JSON(http.StatusOK, game)
	}
}

// Join game
// @Summary Join game
// @Description Join to a battleship game instance
// @Tags Game
// @Accept json
// @Produce json
// @Param request body dto.JoinGameRequest true "Join Game Request"
// @Success 200 {object} dto.JoinGameResponse "Join Game successfully"
// @Router /api/v1/game/join [post]
func (r *GameControllerImpl) JoinGame() func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		request := new(dto.CreateGameRequest)
		if err := ctx.Bind(request); err != nil {
			log.Warn().Err(err).Msg("Bad request")
			return battle_error.BadRequest1(err.Error())
		}

		id, err := r.gameService.CreateGame(*request)
		if err != nil {
			log.Info().Str("userId", request.UserId).Err(err).Msg("cannot join game")
			return err
		}

		return ctx.JSON(http.StatusOK, dto.CreateGameResponse{
			GameId: id,
		})
	}
}
