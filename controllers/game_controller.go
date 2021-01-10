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
	ChangeTurn(ctx echo.Context) error
	RevealEnemyFields(ctx echo.Context) error
	Explode(ctx echo.Context) error
}

type GameControllerImpl struct {
	gameService service.GameService
}

func NewGameControllerImpl(gameService service.GameService) GameControllerImpl {
	return GameControllerImpl{
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
// @Success 200 {object} dto.GetGameResponse "Create Game Response"
// @Router /api/v1/game [post]
func (r GameControllerImpl) CreateGame(ctx echo.Context) error {
	request := new(dto.CreateGameRequest)
	if err := ctx.Bind(request); err != nil {
		log.Warn().Err(err).Msg("Bad request")
		return dto.BadRequest1(err.Error())
	}
	err := request.ValidateAndUnmask()
	if err != nil {
		return err
	}
	res, err := r.gameService.CreateGame(*request)
	if err != nil {
		log.Info().Str("userId", request.UserId).Err(err).Msg("cannot create game")
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
// @Param game_id path string true "Game Id"
// @Param user_id query string true "User Id"
// @Success 200 {object} dto.GetGameResponse "Get Game Response"
// @Router /api/v1/game/{game_id} [get]
func (r GameControllerImpl) GetGame(ctx echo.Context) error {
	gameId := ctx.Param("game_id")
	if gameId == "" {
		log.Warn().Msg("Bad request")
		return dto.BadRequest1("game_id must has value")
	}

	userId := ctx.QueryParam("user_id")
	if userId == "" {
		log.Warn().Msg("Bad request")
		return dto.BadRequest1("user_id must has value")
	}
	request := dto.GetGameRequest{
		UserGameRequest: dto.UserGameRequest{UserId: utils.MaskId(userId), GameId: utils.MaskId(gameId)},
	}
	game, err := r.gameService.GetGame(request)
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
// @Success 200 {object} dto.GetGameResponse "Get Game response"
// @Router /api/v1/game/join [post]
func (r GameControllerImpl) JoinGame(ctx echo.Context) error {
	request := new(dto.JoinGameRequest)
	if err := ctx.Bind(request); err != nil {
		log.Warn().Err(err).Msg("Bad request")
		return dto.BadRequest1(err.Error())
	}
	err := request.ValidateAndUnmask()
	if err != nil {
		return err
	}
	response, err := r.gameService.JoinGame(*request)
	if err != nil {
		log.Info().Str("userId", request.UserId).Err(err).Msg("cannot join game")
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
// @Success 200 {object} dto.SubmitShipsLocationsResponse "Submit Ships Locations Response"
// @Router /api/v1/game/submit-ships [post]
func (r GameControllerImpl) SubmitShipsLocations(ctx echo.Context) error {
	request := new(dto.SubmitShipsLocationsRequest)
	if err := ctx.Bind(request); err != nil {
		log.Warn().Err(err).Msg("Bad request")
		return dto.BadRequest1(err.Error())
	}
	err := request.ValidateAndUnmask()
	if err != nil {
		return err
	}

	response, err := r.gameService.SubmitShipsLocations(*request)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, response)
}

// Move ship
// @Summary Move ship
// @Description Move ship to new location
// @Tags Game
// @Accept json
// @Produce json
// @Param request body dto.MoveShipRequest true "Move Ship Request"
// @Success 200 {object} dto.MoveShipResponse "Move Ship Response"
// @Router /api/v1/game/move-ship [post]
func (r GameControllerImpl) MoveShip(ctx echo.Context) error {

	request := new(dto.MoveShipRequest)
	if err := ctx.Bind(request); err != nil {
		log.Warn().Err(err).Msg("Bad request")
		return dto.BadRequest1(err.Error())
	}
	err := request.ValidateAndUnmask()
	if err != nil {
		return err
	}
	response, err := r.gameService.MoveShip(*request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

// Change turn
// @Summary Change turn
// @Description Change turn
// @Tags Game
// @Accept json
// @Produce json
// @Param request body dto.ChangeTurnRequest true "Change Turn Request"
// @Success 200 {object} dto.ChangeTurnResponse "Change Turn Response"
// @Router /api/v1/game/change-turn [post]
func (r GameControllerImpl) ChangeTurn(ctx echo.Context) error {

	request := new(dto.ChangeTurnRequest)
	if err := ctx.Bind(request); err != nil {
		log.Warn().Err(err).Msg("Bad request")
		return dto.BadRequest1(err.Error())
	}
	err := request.ValidateAndUnmask()
	if err != nil {
		return err
	}
	response, err := r.gameService.ChangeTurn(*request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

// Reveal enemy fields
// @Summary Reveal enemy fields
// @Description Reveal enemy fields
// @Tags Game
// @Accept json
// @Produce json
// @Param request body dto.RevealEnemyFieldsRequest true "Reveal enemy fields request"
// @Success 200 {object} dto.RevealEnemyFieldsResponse "Reveal Enemy Fields Response"
// @Router /api/v1/game/reveal [post]
func (r GameControllerImpl) RevealEnemyFields(ctx echo.Context) error {

	request := new(dto.RevealEnemyFieldsRequest)
	if err := ctx.Bind(request); err != nil {
		log.Warn().Err(err).Msg("Bad request")
		return dto.BadRequest1(err.Error())
	}
	err := request.ValidateAndUnmask()
	if err != nil {
		return err
	}
	response, err := r.gameService.Reveal(*request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}

// Explode a slot
// @Summary Explode a slot
// @Description Explode a slot
// @Tags Game
// @Accept json
// @Produce json
// @Param request body dto.ExplodeRequest true "Explode request"
// @Success 200 {object} dto.ExplodeResponse "Explode Response"
// @Router /api/v1/game/explode [post]
func (r GameControllerImpl) Explode(ctx echo.Context) error {

	request := new(dto.ExplodeRequest)
	if err := ctx.Bind(request); err != nil {
		log.Warn().Err(err).Msg("Bad request")
		return dto.BadRequest1(err.Error())
	}
	err := request.ValidateAndUnmask()
	if err != nil {
		return err
	}
	response, err := r.gameService.Explode(*request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response)
}
