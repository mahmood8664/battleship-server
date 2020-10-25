package service

import (
	"battleship/battle_error"
	"battleship/db/dao"
	"battleship/dto"
	"battleship/model"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
)

type GameService interface {
	CreateGame(request dto.CreateGameRequest) (id string, err error)
	GetGame(gameId string) (game *dto.GameDto, err error)
	JoinGame(request dto.JoinGameRequest) (response *dto.JoinGameResponse, err error)
}

type GameServiceImpl struct {
	gameDao dao.GameDao
	userDao dao.UserDao
}

func NewGameServiceImpl(gameDao dao.GameDao, userDao dao.UserDao) *GameServiceImpl {
	return &GameServiceImpl{
		gameDao: gameDao,
		userDao: userDao,
	}
}

func (r *GameServiceImpl) CreateGame(request dto.CreateGameRequest) (id string, err error) {

	user, err := r.userDao.GetOne(request.UserId)
	if err != nil {
		log.Info().Str("userId", request.UserId).Err(err).Msg("cannot insert user")
		return "", err
	}

	game := model.Game{
		Side1:          &user.Id,
		Status:         model.Init,
		CreateDate:     time.Now(),
		MoveTimeoutSec: 10,
		Turn:           int((rand.Uint32() % 2) + 1),
	}
	return r.gameDao.Insert(&game)
}

func (r *GameServiceImpl) GetGame(gameId string) (game *dto.GameDto, err error) {
	game = new(dto.GameDto)
	g, err := r.gameDao.GetOne(gameId)
	if err == nil {
		game.FromGame(g)
	} else {
		log.Info().Str("gameId", gameId).Err(err).Msg("cannot get Game")
	}
	return game, err
}

func (r *GameServiceImpl) JoinGame(request dto.JoinGameRequest) (response *dto.JoinGameResponse, err error) {

	response = new(dto.JoinGameResponse)

	user, err := r.userDao.GetOne(request.UserId)
	if err != nil {
		log.Info().Str("userId", request.UserId).Str("gameId", request.GameId).Err(err).Msg("cannot get user")
		return response, err
	}

	game, err := r.gameDao.GetOne(request.GameId)
	if err != nil {
		log.Info().Str("userId", request.UserId).Str("gameId", request.GameId).Err(err).Msg("cannot get game")
		return response, err
	}

	if game.Status != model.Init {
		log.Info().Str("userId", request.UserId).Str("gameId", request.GameId).Err(err).Msg("game status is not suitable for joining")
		return response, battle_error.BadRequest2("Game state is not Init", battle_error.InvalidGameStatus)
	}

	if request.UserId != game.Side1.Hex() {
		game.Side2 = &user.Id
		err = r.gameDao.Update(game)
		if err != nil {
			log.Info().Str("userId", request.UserId).Str("gameId", request.GameId).Err(err).Msg("cannot update game")
		}
	}

	response.Game.FromGame(game)
	return response, err
}
