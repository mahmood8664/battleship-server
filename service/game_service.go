package service

import (
	"battleship/cache"
	"battleship/db/dao"
	"battleship/dto"
	"battleship/error_codes"
	"battleship/events/outgoing_events"
	"battleship/model"
	"battleship/utils"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"time"
)

type GameService interface {
	CreateGame(request dto.CreateGameRequest) (response dto.GetGameResponse, err error)
	GetGame(request dto.GetGameRequest) (game dto.GetGameResponse, err error)
	JoinGame(request dto.JoinGameRequest) (response dto.GetGameResponse, err error)
	SubmitShipsLocations(request dto.SubmitShipsLocationsRequest) (response dto.SubmitShipsLocationsResponse, err error)
	ChangeTurn(request dto.ChangeTurnRequest) (response dto.ChangeTurnResponse, err error)
	MoveShip(request dto.MoveShipRequest) (response dto.MoveShipResponse, err error)
	Reveal(request dto.RevealEnemyFieldsRequest) (response dto.RevealEnemyFieldsResponse, err error)
	Explode(request dto.ExplodeRequest) (response dto.ExplodeResponse, err error)
	SocketConnect(event dto.Event, socketConn *websocket.Conn) error
}

type GameServiceImpl struct {
	gameDao      dao.GameDao
	userDao      dao.UserDao
	gameEventDao dao.GameEventDao
	eventHandler outgoing_events.OutgoingEventHandler
}

func NewGameServiceImpl(gameDao dao.GameDao, userDao dao.UserDao, gameEventDao dao.GameEventDao,
	eventHandler outgoing_events.OutgoingEventHandler) GameServiceImpl {
	return GameServiceImpl{
		gameDao:      gameDao,
		userDao:      userDao,
		gameEventDao: gameEventDao,
		eventHandler: eventHandler,
	}
}

func (r GameServiceImpl) CreateGame(request dto.CreateGameRequest) (response dto.GetGameResponse, err error) {
	response = dto.GetGameResponse{}
	user, err := r.userDao.GetOne(request.UserId)
	if err != nil {
		log.Info().Str("userId", request.UserId).Err(err).Msg("cannot insert user")
		return response, err
	}

	rand.Seed(time.Now().UnixNano())

	fields := make(map[int]bool, 100)
	for i := 0; i < 100; i++ {
		fields[i] = true
	}

	game := model.Game{
		Side1User:      &user.Id,
		Side2User:      nil,
		Status:         model.Init,
		LastMoveTime:   time.Now(),
		CreateDate:     time.Now(),
		MoveTimeoutSec: request.MoveTimeout,
		Turn:           int((rand.Uint32() % 2) + 1),
		State: model.GameState{
			Side1Ships:         map[int]bool{},
			Side1Ground:        fields,
			Side2Ships:         map[int]bool{},
			Side2Ground:        fields,
			Side2RevealedShips: map[int]bool{},
			Side1RevealedShips: map[int]bool{},
		},
		WinnerUser: nil,
	}

	gameId, err := r.gameDao.Insert(game)
	if err != nil {
		return response, err
	}

	gameObjectId, err := primitive.ObjectIDFromHex(gameId)
	if err != nil {
		log.Warn().Err(err).Msg("")
		return response, err
	}

	_, err = r.gameEventDao.Insert(model.GameEvent{
		Time:   time.Now(),
		Type:   model.JoinGame,
		GameId: gameObjectId,
		UserId: &user.Id,
	})

	if err != nil {
		log.Warn().Err(err).Msg("")
		return response, err
	}

	gm, err := r.gameDao.GetOne(gameId)
	if err != nil {
		return response, err
	}

	response.Game = new(dto.GameDto)
	response.Game.FromGame(gm, request.UserId)
	response.Ok = true
	return response, nil
}

func (r GameServiceImpl) GetGame(request dto.GetGameRequest) (gameResponse dto.GetGameResponse, err error) {
	gameResponse = dto.GetGameResponse{}

	g, err := r.gameDao.GetOne(request.GameId)
	if err == nil {
		if request.UserId != "" && g.Side1User.Hex() != request.UserId && g.Side2User.Hex() != request.UserId {
			log.Error().Str("user_id", request.UserId).Msg("user does not have access to perform this operation")
			return gameResponse, dto.Forbidden1("cannot perform the operation")
		}

		gameResponse.Game = new(dto.GameDto)
		gameResponse.Game.FromGame(g, request.UserId)
		gameResponse.Ok = true
	} else {
		log.Info().Str("gameId", request.GameId).Err(err).Msg("cannot get Game")
	}
	return gameResponse, err
}

func (r GameServiceImpl) JoinGame(request dto.JoinGameRequest) (response dto.GetGameResponse, err error) {

	response = dto.GetGameResponse{}

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

	if game.Status == model.Init {
		if game.Side1User == nil || game.Side1User.IsZero() {
			log.Error().Str("userId", request.UserId).Str("gameId", request.GameId).
				Msg("user wants to join but side1 has not joint")
			return response, dto.BadRequest2("user wants to join but other side has not joint", error_codes.InvalidGameStatus)
		}

		if game.Side1User.Hex() != request.UserId {
			game.Side2User = &user.Id
			game.Status = model.Joined

			err = r.gameDao.Update(game)
			if err != nil {
				log.Info().Str("userId", request.UserId).Str("gameId", request.GameId).Err(err).
					Msg("cannot update game")
				return response, err
			}

			_, err = r.gameEventDao.Insert(model.GameEvent{
				Time:   time.Now(),
				Type:   model.JoinGame,
				GameId: game.Id,
				UserId: &user.Id,
			})
		}

		response.Game = new(dto.GameDto)
		response.Game.FromGame(game, request.UserId)
		response.Ok = true
		return response, err
	}
	log.Info().Str("userId", request.UserId).Str("gameId", request.GameId).Str("status", string(game.Status)).Err(err).
		Msg("game status is not suitable for joining")
	return response, dto.BadRequest2("game status is not suitable for joining", error_codes.InvalidGameStatus)
}

func (r GameServiceImpl) SubmitShipsLocations(request dto.SubmitShipsLocationsRequest) (response dto.SubmitShipsLocationsResponse, err error) {
	response = dto.SubmitShipsLocationsResponse{}
	game, err := r.gameDao.GetOne(request.GameId)
	if err != nil {
		log.Warn().Err(err).Str("game_id", request.GameId).Msg("error in get game by id")
		return response, err
	}

	if game.Status != model.Joined {
		log.Warn().Msg("game status is not init")
		return response, dto.BadRequest2("Game status is not valid", error_codes.InvalidGameStatus)
	}

	ships := make(map[int]bool)
	for _, element := range request.ShipsIndexes {
		if element < 0 || element >= 100 {
			log.Error().Str("game_id", request.GameId).Str("user_id", request.UserId).
				Int("index", element).Msg("ship indexes must be between 1 and 100")
			return response, dto.BadRequest2("ship indexes must be between 1 and 100", error_codes.InvalidShipIndexValue)
		}
		ships[element] = true
	}

	if len(ships) != 10 {
		log.Error().Str("game_id", request.GameId).Str("user_id", request.UserId).Msg("repeated index is not allowed in ship indexes")
		return response, dto.BadRequest2("repeated index is not allowed in ship indexes", error_codes.InvalidShipIndexValue)
	}

	var otherSide string
	if game.Side1User != nil && game.Side1User.Hex() == request.UserId {
		if len(game.State.Side1Ships) > 0 {
			log.Error().Str("user_id", request.UserId).Str("game_id", request.GameId).
				Msg("user already has chosen his/her ships location")
			return response, dto.Duplicate1("user already has chosen his/her ships location")
		}
		game.State.Side1Ships = ships
		otherSide = game.Side2User.Hex()

		err := r.persistInitialShipLocationEvent(&game.Id, game.Side1User, utils.GetMapKeySlice(ships))
		if err != nil {
			log.Error().Msg("cannot save initial ship location event")
		}
	} else if game.Side2User != nil && game.Side2User.Hex() == request.UserId {
		if len(game.State.Side2Ships) > 0 {
			log.Error().Str("user_id", request.UserId).Str("game_id", request.GameId).
				Msg("user already has chosen his/her ships location")
			return response, dto.Duplicate1("user already has chosen his/her ships location")
		}
		game.State.Side2Ships = ships
		otherSide = game.Side1User.Hex()

		err := r.persistInitialShipLocationEvent(&game.Id, game.Side2User, utils.GetMapKeySlice(ships))
		if err != nil {
			log.Error().Msg("cannot save initial ship location event")
		}
	} else {
		log.Error().Str("user_id", request.UserId).Str("game_id", request.GameId).
			Msg("user is not belong to this game")
		return response, dto.BadRequest1("user is not belong to this game")
	}

	if len(game.State.Side2Ships) == 10 && len(game.State.Side1Ships) == 10 {
		game.Status = model.Start
	}

	game.LastMoveTime = time.Now()

	err = r.gameDao.Update(game)
	if err != nil {
		log.Error().Str("user_id", request.UserId).Str("game_id", request.GameId).
			Err(err).Msg("cannot update game")
		return response, err
	}

	if game.Status == model.Start {
		gameDto := dto.GameDto{}
		gameDto.FromGame(game, otherSide)
		err := r.eventHandler.GameStart(dto.GameStartEvent{
			Game: gameDto,
		})
		if err != nil {
			log.Error().Str("game_id", game.Id.Hex()).Err(err).Msg("cannot send event")
			return response, err
		}
	}

	response.Ok = true
	response.GameStatue = game.Status
	response.Turn = game.Turn
	return response, nil
}

func (r GameServiceImpl) persistInitialShipLocationEvent(gameId *primitive.ObjectID, userId *primitive.ObjectID, initialShipLocations []int) error {
	event := model.GameEvent{
		Type:                  model.InitialShipsLocations,
		InitialShipsLocations: initialShipLocations,
		Time:                  time.Now(),
		UserId:                gameId,
		GameId:                *userId,
	}
	_, err := r.gameEventDao.Insert(event)
	if err != nil {
		log.Error().Str("game_id", gameId.Hex()).Str("user_id", userId.Hex()).
			Msg("cannot save initial ship location event")
		return err
	}
	return nil
}

func (r GameServiceImpl) ChangeTurn(request dto.ChangeTurnRequest) (response dto.ChangeTurnResponse, err error) {
	response = dto.ChangeTurnResponse{}

	game, userId, otherSideUserId, err := r.getGameCheckItWithUserAndChangeTurn(request)
	if err != nil {
		if errors.Is(err, error_codes.NotUserTurn) && !game.Id.IsZero() {
			//in case of leaving game by one of the sides, the present side have to be enable to change turn
			if game.LastMoveTime.Add(time.Duration(game.MoveTimeoutSec+2) * time.Second).Before(time.Now()) {
				if game.Turn == 1 {
					game.Turn = 2
				} else {
					game.Turn = 1
				}

				gameEvents, err := r.gameEventDao.FindMany(request.GameId)
				if err != nil {
					log.Err(err).Str("game_id", request.GameId).Str("user_id", request.UserId).
						Msg("error in getting game events")
					return response, err
				}
				//if four last moves belong to one side, other side is not present anymore and game will be finished
				if len(gameEvents) >= 4 {
					if gameEvents[0].UserId.Hex() == gameEvents[1].UserId.Hex() &&
						gameEvents[0].UserId.Hex() == gameEvents[2].UserId.Hex() &&
						gameEvents[0].UserId.Hex() == gameEvents[3].UserId.Hex() &&
						gameEvents[0].UserId.Hex() == gameEvents[4].UserId.Hex() {
						game.Status = model.Finished
						game.WinnerUser = gameEvents[0].UserId

						err := r.eventHandler.EndGame(dto.EndGameEvent{
							GameId:       utils.MaskId(request.GameId),
							WinnerUserId: utils.MaskId(game.WinnerUser.Hex()),
						})
						if err != nil {
							log.Error().Str("game_id", request.GameId).Msg("cannot send end game event")
						}
					}
				}
			} else {
				log.Err(err).Str("game_id", request.GameId).Str("user_id", request.UserId).
					Msg("it's not user turn")
				return response, err
			}
		} else {
			log.Error().Err(err).Msg("error in checking game")
			return response, err
		}
	}

	err = r.gameDao.Update(game)
	if err != nil {
		log.Err(err).Str("game_id", request.GameId).Str("user_id", request.UserId).
			Msg("error in updating game")
	}

	if game.Status != model.Finished {
		_, err = r.gameEventDao.Insert(model.GameEvent{
			Time:   time.Now(),
			Type:   model.ChangeTurn,
			GameId: game.Id,
			UserId: &userId,
		})

		err = r.eventHandler.ChangeTurn(dto.GameChangeTurnEvent{
			GameId: utils.MaskId(request.GameId),
			UserId: utils.MaskId(otherSideUserId.Hex()),
		})
		if err != nil {
			log.Err(err).Msg("cannot send change turn event")
		}
	}
	response.Ok = true
	return response, nil
}

func (r GameServiceImpl) MoveShip(request dto.MoveShipRequest) (response dto.MoveShipResponse, err error) {
	response = dto.MoveShipResponse{}
	game, userId, otherSide, err := r.getGameCheckItWithUserAndChangeTurn(request)
	if err != nil {
		log.Error().Err(err).Msg("error in checking game")
		return response, err
	}

	if game.Side1User.Hex() == request.UserId {
		err := game.MoveShipSide1(request.OldShipIndex, request.NewShipIndex)
		if err != nil {
			log.Err(err).Str("game_id", request.GameId).Str("user_id", request.UserId).
				Err(err).Msg("error in move ship in side 1")
			return response, err
		}
	} else if game.Side2User.Hex() == request.UserId {
		err := game.MoveShipSide2(request.OldShipIndex, request.NewShipIndex)
		if err != nil {
			log.Err(err).Str("game_id", request.GameId).Str("user_id", request.UserId).
				Err(err).Msg("error in move ship in side 2")
			return response, err
		}
	}

	err = r.gameDao.Update(game)
	if err != nil {
		log.Error().Str("game_id", request.GameId).Msg("cannot update game state")
		return response, err
	}

	err = r.submitMoveShipEvent(game.Id, userId, request.OldShipIndex, request.NewShipIndex)
	if err != nil {
		log.Error().Str("game_id", request.GameId).Str("user_id", request.UserId).
			Msg("cannot save move ship event")
		return response, err
	}

	err = r.eventHandler.MoveShip(dto.ShipMovedEvent{
		GameId:       utils.MaskId(request.GameId),
		UserId:       utils.MaskId(otherSide.Hex()),
		OldShipIndex: request.OldShipIndex,
	})
	if err != nil {
		log.Error().Str("game_id", request.GameId).Str("user_id", request.UserId).
			Msg("cannot send ship move event")
		return response, err
	}

	response.Ok = true
	return response, nil
}

func (r GameServiceImpl) submitMoveShipEvent(gameId primitive.ObjectID, userId primitive.ObjectID, oldShipIndex int,
	newShipIndex int) error {
	event := model.GameEvent{
		Type:         model.MoveShip,
		MoveShipFrom: &oldShipIndex,
		MoveShipTo:   &newShipIndex,
		Time:         time.Now(),
		UserId:       &userId,
		GameId:       gameId,
	}
	_, err := r.gameEventDao.Insert(event)
	if err != nil {
		log.Error().Str("game_id", gameId.Hex()).Str("user_id", userId.Hex()).
			Msg("cannot save move ship event")
		return err
	}
	return nil
}

func (r GameServiceImpl) Reveal(request dto.RevealEnemyFieldsRequest) (response dto.RevealEnemyFieldsResponse, err error) {
	response = dto.RevealEnemyFieldsResponse{}
	game, userId, otherSide, err := r.getGameCheckItWithUserAndChangeTurn(request)
	if err != nil {
		log.Error().Err(err).Msg("error in checking game")
		return response, err
	}

	var revealedShipsIndexes []int
	if game.Side1User.Hex() == request.UserId {
		revealedShipsIndexes = game.RevealSlotSide2(request.Index)
	} else if game.Side2User.Hex() == request.UserId {
		revealedShipsIndexes = game.RevealSlotSide1(request.Index)
	}

	err = r.gameDao.Update(game)
	if err != nil {
		log.Warn().Str("game_id", request.GameId).Msg("cannot update game")
		return response, err
	}

	err = r.PersistRevealEvent(game.Id, userId, request.Index, revealedShipsIndexes)
	if err != nil {
		log.Error().Msg("cannot save reveal event")
	}

	err = r.eventHandler.Reveal(dto.RevealEvent{
		UserId:        utils.MaskId(otherSide.Hex()),
		GameId:        utils.MaskId(request.GameId),
		RevealedShips: revealedShipsIndexes,
		Slots:         model.FindNeighborIndexes(request.Index),
	})
	if err != nil {
		log.Error().Str("game_id", request.GameId).Str("user_id", request.UserId).
			Msg("cannot send reveal event")
	}

	response.Ok = true
	response.RevealedShipIndexes = revealedShipsIndexes
	return response, nil
}

func (r GameServiceImpl) PersistRevealEvent(gameId primitive.ObjectID, userId primitive.ObjectID, index int, revealedShipIndexes []int) error {
	event := model.GameEvent{
		Type:               model.Reveal,
		Time:               time.Now(),
		DiscoverEnemy:      model.FindNeighborIndexes(index),
		DiscoverEnemyShips: revealedShipIndexes,
		UserId:             &userId,
		GameId:             gameId,
	}
	_, err := r.gameEventDao.Insert(event)
	if err != nil {
		log.Error().Str("game_id", gameId.Hex()).Str("user_id", userId.Hex()).
			Msg("cannot save reveal event")
		return err
	}
	return nil
}

func (r GameServiceImpl) Explode(request dto.ExplodeRequest) (response dto.ExplodeResponse, err error) {
	response = dto.ExplodeResponse{}

	game, userId, otherSide, err := r.getGameCheckItWithUserAndChangeTurn(request)
	if err != nil {
		log.Error().Err(err).Msg("error in checking game")
		return response, err
	}

	if game.Side1User.Hex() == request.UserId {
		response.HasShip = game.ExplodeSide2(request.Index)
	} else if game.Side2User.Hex() == request.UserId {
		response.HasShip = game.ExplodeSide1(request.Index)
	}

	if response.HasShip {
		if game.Turn == 1 {
			game.Turn = 2
		} else {
			game.Turn = 1
		}
	}

	err = r.PersistExplosionEvent(game.Id, userId, request.Index, !response.HasShip)
	if err != nil {
		log.Error().Msg("cannot save explosion event")
	}

	err = r.gameDao.Update(game)
	if err != nil {
		log.Warn().Str("game_id", request.GameId).Msg("cannot update game")
		return response, err
	}

	err = r.eventHandler.Explosion(dto.ExplosionEvent{
		GameId: utils.MaskId(request.GameId),
		UserId: utils.MaskId(otherSide.Hex()),
		Index:  request.Index,
	})

	if game.Status == model.Finished && game.WinnerUser.IsZero() == false {
		err := r.eventHandler.EndGame(dto.EndGameEvent{
			GameId:       utils.MaskId(request.GameId),
			WinnerUserId: utils.MaskId(game.WinnerUser.Hex()),
		})
		if err != nil {
			log.Error().Str("game_id", request.GameId).Msg("cannot send end game event")
		}
	}

	if err != nil {
		log.Error().Str("game_id", request.GameId).Str("user_id", request.UserId).
			Msg("cannot send explosion event")
	}

	response.Ok = true
	return response, nil
}

func (r GameServiceImpl) PersistExplosionEvent(gameId primitive.ObjectID, userId primitive.ObjectID, index int, empty bool) error {
	var event model.GameEvent
	if empty {
		event = model.GameEvent{
			Type:           model.EmptyExplosion,
			Time:           time.Now(),
			UserId:         &userId,
			GameId:         gameId,
			EmptyExplosion: &index,
		}
	} else {
		event = model.GameEvent{
			Type:      model.Explosion,
			Time:      time.Now(),
			UserId:    &userId,
			GameId:    gameId,
			Explosion: &index,
		}
	}
	_, err := r.gameEventDao.Insert(event)
	if err != nil {
		log.Error().Str("game_id", gameId.Hex()).Str("user_id", userId.Hex()).
			Msg("cannot save explosion event")
		return err
	}
	return nil
}

func (r GameServiceImpl) SocketConnect(event dto.Event, socketConn *websocket.Conn) error {
	request := new(dto.UserConnectEvent)
	err := json.Unmarshal([]byte(event.Payload), request)
	if err != nil {
		return dto.ParseError(err)
	}
	err = request.ValidateAndUnmask()
	if err != nil {
		return err
	}

	game, err := r.gameDao.GetOne(request.GameId)
	if err != nil {
		return err
	}

	if _, ok := cache.GameCache.Cache[request.GameId]; ok == false {
		cache.GameCache.Cache[request.GameId] = cache.GameData{}
	}

	gameData := cache.GameCache.Cache[request.GameId]

	if game.Side1User != nil && request.UserId == game.Side1User.Hex() {
		gameData.Side1Socket = socketConn
		gameData.Side1UserId = request.UserId
		if game.Status == model.Joined {
			err := r.eventHandler.GameConnect(dto.GameConnect{
				GameId: utils.MaskId(game.Id.Hex()),
				UserId: utils.MaskId(game.Side2User.Hex()),
			})
			if err != nil {
				log.Error().Str("user_id", request.UserId).Str("game_id", request.GameId).
					Msg("cannot send user connect event")
			}
		}
	} else if game.Side2User != nil && request.UserId == game.Side2User.Hex() {
		gameData.Side2Socket = socketConn
		gameData.Side2UserId = request.UserId
		if game.Status == model.Joined {
			err := r.eventHandler.GameConnect(dto.GameConnect{
				GameId: utils.MaskId(game.Id.Hex()),
				UserId: utils.MaskId(game.Side1User.Hex()),
			})
			if err != nil {
				log.Error().Str("user_id", request.UserId).Str("game_id", request.GameId).
					Msg("cannot send user connect event")
			}
		}
	} else {
		log.Error().Str("user_id", request.UserId).Str("game_id", request.GameId).Msg("user does not belong to game")
		return dto.BadRequest1("user does not belong to game")
	}
	cache.GameCache.Cache[game.Id.Hex()] = gameData
	log.Debug().Str("game_id", request.GameId).Str("user_id", request.UserId).Msg("create socket successfully")
	return nil
}

func (r GameServiceImpl) getGameCheckItWithUserAndChangeTurn(request dto.UserGame) (game model.Game, userId primitive.ObjectID, otherSide primitive.ObjectID, err error) {
	game, err = r.gameDao.GetOne(request.GetGameId())
	if err != nil {
		log.Error().Str("game_id", request.GetGameId()).Str("user_id", request.GetUserId()).
			Msg("cannot find game")
		return game, userId, otherSide, err
	}

	if game.Status != model.Start {
		log.Error().Str("game_id", request.GetGameId()).Msg("game is not in start status")
		return game, userId, otherSide, dto.BadRequest2("game is not started", error_codes.InvalidGameStatus)
	}

	if game.LastMoveTime.Add(1 * time.Minute).Before(time.Now()) {
		log.Error().Str("game_id", request.GetGameId()).Msg("game is already finished")
		return game, userId, otherSide, dto.BadRequest2("game is finished", error_codes.GameIsFinished)
	}

	if game.Side1User.Hex() == request.GetUserId() {
		otherSide = *game.Side2User
		userId = *game.Side1User
		if game.Turn == 1 {
			game.Turn = 2
		} else {
			return game, userId, otherSide, error_codes.NotUserTurn
		}

	} else if game.Side2User.Hex() == request.GetUserId() {
		otherSide = *game.Side1User
		userId = *game.Side2User
		if game.Turn == 2 {
			game.Turn = 1
		} else {
			return game, userId, otherSide, error_codes.NotUserTurn
		}
	} else {
		log.Error().Str("game_id", request.GetGameId()).Str("user_id", request.GetUserId()).
			Msg("user does not belong to this game")
		return game, userId, otherSide, dto.Forbidden1("user does not belong to this game")
	}
	game.LastMoveTime = time.Now()
	return game, userId, otherSide, nil
}
