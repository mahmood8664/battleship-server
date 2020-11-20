package events

import (
	"battleship/cache"
	"battleship/dto"
	"battleship/service"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type ConnectionEventHandler interface {
	UserConnect(event dto.Event, socketConn *websocket.Conn) error
}

type ConnectionEventHandlerImpl struct {
	gameService service.GameService
}

func NewConnectionEventHandlerImpl(gameService service.GameService) *ConnectionEventHandlerImpl {
	return &ConnectionEventHandlerImpl{
		gameService: gameService,
	}
}

func (r *ConnectionEventHandlerImpl) UserConnect(event dto.Event, socketConn *websocket.Conn) error {
	request := new(dto.SocketConnect)
	err := json.Unmarshal([]byte(event.Payload), request)
	if err != nil {
		return dto.ParseError(err)
	}

	gameResponse, err := r.gameService.GetGame(request.GameId)
	if err != nil {
		return err
	}

	if _, ok := cache.GameCache.Cache[request.GameId]; ok == false {
		cache.GameCache.Cache[request.GameId] = cache.GameData{}
	}

	gameData := cache.GameCache.Cache[request.GameId]

	if gameResponse.Game.Side1UserId != nil && request.UserId == *gameResponse.Game.Side1UserId {
		gameData.Side1Socket = socketConn
		gameData.Side1UserId = request.UserId
		sendConnectEvent(request, gameData.Side2Socket)
	} else if gameResponse.Game.Side2UserId != nil && request.UserId == *gameResponse.Game.Side2UserId {
		gameData.Side2Socket = socketConn
		gameData.Side2UserId = request.UserId
		sendConnectEvent(request, gameData.Side1Socket)
	} else {
		log.Error().Str("user_id", request.UserId).
			Str("game_id", request.GameId).
			Msg("user does not belong to game")
		return dto.BadRequest1("user does not belong to game")
	}
	cache.GameCache.Cache[request.GameId] = gameData
	log.Debug().Str("game_id", request.GameId).Str("user_id", request.UserId).Msg("create socket successfully")
	return nil
}

func sendConnectEvent(request *dto.SocketConnect, socket *websocket.Conn) {
	if socket != nil {

		socketConnectEventBytes, err := json.Marshal(request)
		if err != nil {
			log.Info().Err(err).Str("game_id", request.GameId).Str("user_id", request.UserId).
				Msg("error in marshalling Socket Connect message")
		}

		event := dto.Event{
			Type:    dto.Connect,
			Payload: string(socketConnectEventBytes),
		}

		eventBytes, err := json.Marshal(event)
		if err != nil {
			log.Info().Err(err).Str("game_id", request.GameId).Str("user_id", request.UserId).
				Msg("error in marshalling event message")
		}

		err = socket.WriteMessage(websocket.TextMessage, eventBytes)
		if err != nil {
			log.Info().Err(err).Str("game_id", request.GameId).Str("user_id", request.UserId).
				Msg("error in sending Socket Connect message")
		}

	}
}
