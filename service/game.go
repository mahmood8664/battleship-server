package service

import (
	"battleship/db/dao"
	"battleship/dto"
	"battleship/error_codes"
	"battleship/events/outgoing_events"
	"battleship/model"
	"battleship/utils"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
)

type GameService interface {
	CreateGame(request dto.CreateGameRequest) (response dto.CreateGameResponse, err error)
	GetGame(gameId string) (game *dto.GetGameResponse, err error)
	JoinGame(request dto.JoinGameRequest) (response *dto.GetGameResponse, err error)
	SubmitShipsLocations(request dto.SubmitShipsLocationsRequest) (response *dto.GetGameResponse, err error)
}

type GameServiceImpl struct {
	gameDao      dao.GameDao
	userDao      dao.UserDao
	eventHandler outgoing_events.OutgoingEventHandler
}

func NewGameServiceImpl(gameDao dao.GameDao, userDao dao.UserDao, eventHandler outgoing_events.OutgoingEventHandler) *GameServiceImpl {
	return &GameServiceImpl{
		gameDao:      gameDao,
		userDao:      userDao,
		eventHandler: eventHandler,
	}
}

func (r *GameServiceImpl) CreateGame(request dto.CreateGameRequest) (response dto.CreateGameResponse, err error) {
	user, err := r.userDao.GetOne(request.GetUnMaskedUserId())
	if err != nil {
		log.Info().Str("userId", request.GetUnMaskedUserId()).Err(err).Msg("cannot insert user")
		return response, err
	}

	rand.Seed(time.Now().UnixNano())

	fields := make(map[int]bool, 200)
	for i := 0; i < 200; i++ {
		fields[i] = true
	}

	game := model.Game{
		Side1:          &user.Id,
		Side2:          nil,
		Status:         model.Init,
		CreateDate:     time.Now(),
		MoveTimeoutSec: 10,
		Turn:           int((rand.Uint32() % 2) + 1),
		State: model.GameState{
			Side1Ships: nil,
			Side1:      fields,
			Side2Ships: nil,
			Side2:      fields,
		},
		Winner:       nil,
		LastMoveTime: nil,
	}
	id, err := r.gameDao.Insert(&game)
	if err != nil {
		return response, err
	}
	response.GameId = utils.MaskId(id)
	response.Ok = true
	return response, nil
}

func (r *GameServiceImpl) GetGame(gameId string) (gameResponse *dto.GetGameResponse, err error) {
	gameResponse = new(dto.GetGameResponse)
	g, err := r.gameDao.GetOne(gameId)
	if err == nil {
		gameResponse.Game = new(dto.GameDto)
		gameResponse.Game.FromGame(g)
	} else {
		log.Info().Str("gameId", gameId).Err(err).Msg("cannot get Game")
	}
	return gameResponse, err
}

func (r *GameServiceImpl) JoinGame(request dto.JoinGameRequest) (response *dto.GetGameResponse, err error) {

	response = new(dto.GetGameResponse)

	user, err := r.userDao.GetOne(request.GetUnMaskedUserId())
	if err != nil {
		log.Info().Str("userId", request.GetUnMaskedUserId()).Str("gameId", request.GetUnMaskedGameId()).Err(err).Msg("cannot get user")
		return response, err
	}

	game, err := r.gameDao.GetOne(request.GetUnMaskedGameId())
	if err != nil {
		log.Info().Str("userId", request.GetUnMaskedUserId()).Str("gameId", request.GetUnMaskedGameId()).Err(err).Msg("cannot get game")
		return response, err
	}

	if game.Status != model.Init {
		log.Info().Str("userId", request.GetUnMaskedUserId()).Str("gameId", request.GetUnMaskedGameId()).Err(err).
			Msg("game status is not suitable for joining")
		return response, dto.BadRequest2("Game state is not Init", error_codes.InvalidGameStatus)
	}

	if request.GetUnMaskedUserId() == game.Side1.Hex() {
		log.Info().Str("userId", request.GetUnMaskedUserId()).Str("gameId", request.GetUnMaskedGameId()).Msg("user already joint")
		//return response, dto.Duplicate0()
	}

	if game.Side2 != nil && request.GetUnMaskedUserId() == game.Side2.Hex() {
		log.Info().Str("userId", request.GetUnMaskedUserId()).Str("gameId", request.GetUnMaskedGameId()).Msg("user already joint")
		//return response, dto.Duplicate0()
	}

	game.Side2 = &user.Id
	err = r.gameDao.Update(game)
	if err != nil {
		log.Info().Str("userId", request.GetUnMaskedUserId()).Str("gameId", request.GetUnMaskedGameId()).Err(err).Msg("cannot update game")
		return response, err
	}

	response.Game = new(dto.GameDto)
	response.Game.FromGame(game)
	response.Ok = true
	return response, err
}

func (r *GameServiceImpl) SubmitShipsLocations(request dto.SubmitShipsLocationsRequest) (response *dto.GetGameResponse, err error) {
	response = new(dto.GetGameResponse)
	game, err := r.gameDao.GetOne(request.GetUnMaskedGameId())
	if err != nil {
		log.Warn().Str("game_id", request.GetUnMaskedGameId()).Msg("error in get game by id")
		return response, err
	}

	if game.Status != model.Init {
		log.Warn().Msg("game status is not init")
		return response, dto.BadRequest2("invalid game status", error_codes.InvalidGameStatus)
	}

	ships := make(map[int]bool)
	for _, element := range request.ShipsIndexes {
		if element < 1 || element > 200 {
			log.Error().Str("game_id", request.GameId).Str("user_id", request.UserId).Msg("ship indexes must be between 1 and 200")
			return response, dto.BadRequest2("ship indexes must be between 1 and 200", error_codes.InvalidShipIndexValue)
		}
		ships[element] = true
	}

	if len(ships) != 10 {
		log.Error().Str("game_id", request.GameId).Str("user_id", request.UserId).Msg("repeated index is not allowed in ship indexes")
		return response, dto.BadRequest2("repeated index is not allowed in ship indexes", error_codes.InvalidShipIndexValue)
	}

	if game.Side1 != nil && game.Side1.Hex() == request.GetUnMaskedUserId() {
		if game.State.Side1Ships != nil {
			log.Error().Str("user_id", request.GetUnMaskedUserId()).Str("game_id", request.GetUnMaskedGameId()).
				Msg("user already has chosen his/her ships location")
			return response, dto.Duplicate1("user already has chosen his/her ships location")
		}
		game.State.Side1Ships = ships

	} else if game.Side2 != nil && game.Side2.Hex() == request.GetUnMaskedUserId() {
		if game.State.Side2Ships != nil {
			log.Error().Str("user_id", request.GetUnMaskedUserId()).Str("game_id", request.GetUnMaskedGameId()).
				Msg("user already has chosen his/her ships location")
			return response, dto.Duplicate0()
		}
		game.State.Side1Ships = ships
	} else {
		log.Error().Str("user_id", request.GetUnMaskedUserId()).Str("game_id", request.GetUnMaskedGameId()).
			Msg("user is not belong to this game")
		return response, dto.BadRequest0()
	}

	if len(game.State.Side2Ships) == 10 && len(game.State.Side1Ships) == 10 {
		game.Status = model.Start
	}

	err = r.gameDao.Update(game)
	if err != nil {
		log.Error().Str("user_id", request.GetUnMaskedUserId()).Str("game_id", request.GetUnMaskedGameId()).
			Msg("cannot update game")
		return response, err
	}

	if len(game.State.Side2Ships) == 10 && len(game.State.Side1Ships) == 10 {

	}

	response.Game = new(dto.GameDto)
	response.Game.FromGame(game)
	response.Ok = true
	return response, nil
}
