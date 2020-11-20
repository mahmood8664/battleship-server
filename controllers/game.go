package controllers

import (
	"battleship/dto"
	"battleship/service"
	"battleship/utils"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

type GameController interface {
	CreateGame(ctx echo.Context) error
	GetGame(ctx echo.Context) error
	JoinGame(ctx echo.Context) error
	MoveShip(ctx echo.Context) error
	SubmitShipsLocations(ctx echo.Context) error
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
func (r *GameControllerImpl) CreateGame(ctx echo.Context) error {
	request := new(dto.CreateGameRequest)
	if err := ctx.Bind(request); err != nil {
		log.Warn().Err(err).Msg("Bad request")
		return dto.BadRequest1(err.Error())
	}

	res, err := r.gameService.CreateGame(*request)
	if err != nil {
		log.Info().Str("userId", request.GetUnMaskedUserId()).Err(err).Msg("cannot create game")
		return err
	}
	return ctx.JSON(http.StatusOK, res)
}

// Get game info
// @Summary Get game
// @Description Get game info
// @Tags Game
// @Accept json
// @Produce json
// @Param game_id path string true " "
// @Success 200 {object} dto.GetGameResponse "Game successfully created"
// @Router /api/v1/game/{game_id} [get]
func (r *GameControllerImpl) GetGame(ctx echo.Context) error {
	gameId := ctx.Param("game_id")
	if gameId == "" {
		log.Warn().Msg("Bad request")
		return dto.BadRequest1("game_id must has value")
	}

	game, err := r.gameService.GetGame(utils.MaskId(gameId))
	if err != nil {
		log.Info().Str("gameId", utils.MaskId(gameId)).Err(err).Msg("cannot get Game")
		return err
	}
	return ctx.JSON(http.StatusOK, game)
}

// Join game
// @Summary Join game
// @Description Join to a battleship game instance
// @Tags Game
// @Accept json
// @Produce json
// @Param request body dto.JoinGameRequest true "Join Game Request"
// @Success 200 {object} dto.GetGameResponse "Join Game successfully"
// @Router /api/v1/game/join [post]
func (r *GameControllerImpl) JoinGame(ctx echo.Context) error {
	request := new(dto.JoinGameRequest)
	if err := ctx.Bind(request); err != nil {
		log.Warn().Err(err).Msg("Bad request")
		return dto.BadRequest1(err.Error())
	}

	response, err := r.gameService.JoinGame(*request)
	if err != nil {
		log.Info().Str("userId", request.GetUnMaskedUserId()).Err(err).Msg("cannot join game")
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

// submit ship locations
// @Summary submit ship locations
// @Description submit ship locations
// @Tags Game
// @Accept json
// @Produce json
// @Param request body dto.SubmitShipsLocationsRequest true "Submit ships Request"
// @Success 200 {object} dto.GetGameResponse "Submit ships successfully"
// @Router /api/v1/game/submit-ships [post]
func (r *GameControllerImpl) SubmitShipsLocations(ctx echo.Context) error {
	request := new(dto.SubmitShipsLocationsRequest)
	if err := ctx.Bind(request); err != nil {
		log.Warn().Err(err).Msg("Bad request")
		return dto.BadRequest1(err.Error())
	}

	if len(request.ShipsIndexes) != 10 {
		log.Error().Msg("ship index size must be 10")
		return dto.BadRequest0()
	}

	response, err := r.gameService.SubmitShipsLocations(*request)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, response)
}

// Move ship
// @Summary Move ship
// @Description Move ship to new locaiton
// @Tags Game
// @Accept json
// @Produce json
// @Param request body dto.MoveShipRequest true "Join Game Request"
// @Success 200 {object} dto.GetGameResponse "Join Game successfully"
// @Router /api/v1/game/move-ship [post]
func (r *GameControllerImpl) MoveShip(ctx echo.Context) error {

	return ctx.JSON(http.StatusOK, "")
}
