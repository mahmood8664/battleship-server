package socket

import (
	"battleship/dto"
	"battleship/events/incoming_events"
	"battleship/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

//goland:noinspection GoNameStartsWithPackageName
type SocketHandler interface {
	CreateSocket(c echo.Context) error
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		fmt.Printf("%+v", r)
		return true
	},
}

//goland:noinspection GoNameStartsWithPackageName
type SocketHandlerImpl struct {
	incomingEventHandler incoming_events.IncomingEventHandler
	gameService          service.GameService
}

func NewSocketHandlerImpl(incomingEventHandler incoming_events.IncomingEventHandler, gameService service.GameService) SocketHandlerImpl {
	return SocketHandlerImpl{
		incomingEventHandler: incomingEventHandler,
		gameService:          gameService,
	}
}

func (r SocketHandlerImpl) CreateSocket(c echo.Context) error {
	return r.connect(c)
}

func (r SocketHandlerImpl) connect(c echo.Context) error {
	gameId := c.QueryParam("game_id")
	userId := c.QueryParam("user_id")
	if gameId == "" || userId == "" {
		log.Error().Str("game_id", gameId).Str("user_id", userId).Msg("game_id or user_id is null in create socket")
		return dto.BadRequest0()
	}
	socketConn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Error().Msg("error in upgrading:" + err.Error())
		return err
	}

	request := new(dto.UserConnectEvent)
	request.GameId = gameId
	request.UserId = userId
	marshal, err := json.Marshal(request)
	if err != nil {
		log.Error().Err(err).Msg("error in marshaling UserConnectEvent")
		return err
	}

	event := dto.Event{
		Type:    dto.Connect,
		Payload: string(marshal),
	}
	err = r.gameService.SocketConnect(event, socketConn)
	if err != nil {
		log.Error().Err(err).Msg("error in NewConnectionHandler")
		return err
	}

	for {
		_, message, err := socketConn.ReadMessage()
		if err != nil {
			log.Warn().Msg("error in reading message:" + err.Error())
			break
		}
		log.Debug().Msg("recv: " + string(message))
		event := new(dto.Event)
		err = json.Unmarshal(message, event)
		if err != nil {
			log.Error().Str("message", string(message)).Msg("message format is not Event")
			//_ = c.Close()
		}

		err = r.incomingEventHandler.HandleEvent(*event)
		if err != nil {
			_ = socketConn.Close()
		}
	}
	return nil
}
