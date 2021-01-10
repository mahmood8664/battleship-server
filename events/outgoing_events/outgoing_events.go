package outgoing_events

import (
	"battleship/cache"
	"battleship/dto"
	"battleship/utils"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type OutgoingEventHandler interface {
	GameConnect(gameConnectEvent dto.GameConnect) error
	GameStart(gameStartEvent dto.GameStartEvent) error
	ChangeTurn(changeTurn dto.GameChangeTurnEvent) error
	MoveShip(shipMovedEvent dto.ShipMovedEvent) error
	Reveal(revealEvent dto.RevealEvent) error
	Explosion(explosionEvent dto.ExplosionEvent) error
	EndGame(endGameEvent dto.EndGameEvent) error
}

type OutgoingEventHandlerImpl struct {
}

func NewOutgoingEventHandlerImpl() OutgoingEventHandlerImpl {
	return OutgoingEventHandlerImpl{}
}

func (r OutgoingEventHandlerImpl) GameConnect(gameConnectEvent dto.GameConnect) error {
	if gameData, ok := cache.GameCache.Cache[utils.MaskId(gameConnectEvent.GameId)]; ok {
		eventBytes, err := dto.MarshalEvent(gameConnectEvent, dto.Connect)
		if err != nil {
			log.Error().Err(err).Msg("cannot marshal GameStartEvent")
			return err
		}

		if gameData.Side1UserId == utils.MaskId(gameConnectEvent.UserId) {
			err = gameData.Side1Socket.WriteMessage(websocket.TextMessage, eventBytes)
		} else if gameData.Side2UserId == utils.MaskId(gameConnectEvent.UserId) {
			err = gameData.Side2Socket.WriteMessage(websocket.TextMessage, eventBytes)
		} else {
			return dto.Forbidden1("user does not belong to game!")
		}

		if err != nil {
			log.Error().Err(err).Msg("cannot send GameStartEvent")
			return err
		}
	}
	return nil
}

func (r OutgoingEventHandlerImpl) GameStart(gameStartEvent dto.GameStartEvent) error {
	if gameData, ok := cache.GameCache.Cache[utils.MaskId(gameStartEvent.Game.Id)]; ok {
		eventBytes, err := dto.MarshalEvent(gameStartEvent, dto.GameStart)
		if err != nil {
			log.Error().Err(err).Msg("cannot marshal GameStartEvent")
			return err
		}

		if gameData.Side1UserId == utils.MaskId(gameStartEvent.Game.UserId) {
			err = gameData.Side1Socket.WriteMessage(websocket.TextMessage, eventBytes)
		} else if gameData.Side2UserId == utils.MaskId(gameStartEvent.Game.UserId) {
			err = gameData.Side2Socket.WriteMessage(websocket.TextMessage, eventBytes)
		} else {
			return dto.Forbidden1("user does not belong to game!")
		}

		if err != nil {
			log.Error().Err(err).Msg("cannot send GameStartEvent")
			return err
		}
	}
	return nil
}

func (r OutgoingEventHandlerImpl) ChangeTurn(changeTurnEvent dto.GameChangeTurnEvent) error {
	if gameData, ok := cache.GameCache.Cache[utils.MaskId(changeTurnEvent.GameId)]; ok {

		eventBytes, err := dto.MarshalEvent(changeTurnEvent, dto.ChangeTurn)
		if err != nil {
			log.Error().Err(err).Msg("cannot marshal GameChangeTurnEvent")
			return err
		}

		if gameData.Side1UserId == utils.MaskId(changeTurnEvent.UserId) {
			err = gameData.Side1Socket.WriteMessage(websocket.TextMessage, eventBytes)
		} else if gameData.Side2UserId == utils.MaskId(changeTurnEvent.UserId) {
			err = gameData.Side2Socket.WriteMessage(websocket.TextMessage, eventBytes)
		} else {
			log.Err(err).Msg("user does not belong to game!")
			return dto.Forbidden1("user does not belong to game!")
		}

		if err != nil {
			log.Err(err).Msg("cannot send GameChangeTurnEvent")
			return err
		}
	}
	return nil
}

func (r OutgoingEventHandlerImpl) MoveShip(shipMovedEvent dto.ShipMovedEvent) error {
	if gameData, ok := cache.GameCache.Cache[utils.MaskId(shipMovedEvent.GameId)]; ok {

		eventBytes, err := dto.MarshalEvent(shipMovedEvent, dto.ShipMoved)
		if err != nil {
			log.Error().Err(err).Msg("cannot marshal ShipMovedEvent")
			return err
		}

		if gameData.Side1UserId == utils.MaskId(shipMovedEvent.UserId) {
			err = gameData.Side1Socket.WriteMessage(websocket.TextMessage, eventBytes)
		} else if gameData.Side2UserId == utils.MaskId(shipMovedEvent.UserId) {
			err = gameData.Side2Socket.WriteMessage(websocket.TextMessage, eventBytes)
		} else {
			return dto.Forbidden1("user does not belong to game!")
		}

		if err != nil {
			log.Err(err).Msg("cannot send ShipMovedEvent")
		}
	}
	return nil
}

func (r OutgoingEventHandlerImpl) Reveal(revealEvent dto.RevealEvent) error {
	if gameData, ok := cache.GameCache.Cache[utils.MaskId(revealEvent.GameId)]; ok {

		eventBytes, err := dto.MarshalEvent(revealEvent, dto.Reveal)
		if err != nil {
			log.Error().Err(err).Msg("cannot marshal RevealEvent")
			return err
		}

		if gameData.Side1UserId == utils.MaskId(revealEvent.UserId) {
			err = gameData.Side1Socket.WriteMessage(websocket.TextMessage, eventBytes)
		} else if gameData.Side2UserId == utils.MaskId(revealEvent.UserId) {
			err = gameData.Side2Socket.WriteMessage(websocket.TextMessage, eventBytes)
		} else {
			return dto.Forbidden1("user does not belong to game!")
		}

		if err != nil {
			log.Err(err).Msg("cannot send RevealEvent")
		}
	}
	return nil
}

func (r OutgoingEventHandlerImpl) Explosion(explosionEvent dto.ExplosionEvent) error {
	if gameData, ok := cache.GameCache.Cache[utils.MaskId(explosionEvent.GameId)]; ok {

		eventBytes, err := dto.MarshalEvent(explosionEvent, dto.Explosion)
		if err != nil {
			log.Error().Err(err).Msg("cannot marshal ExplosionEvent")
			return err
		}

		if gameData.Side1UserId == utils.MaskId(explosionEvent.UserId) {
			err = gameData.Side1Socket.WriteMessage(websocket.TextMessage, eventBytes)
		} else if gameData.Side2UserId == utils.MaskId(explosionEvent.UserId) {
			err = gameData.Side2Socket.WriteMessage(websocket.TextMessage, eventBytes)
		} else {
			return dto.Forbidden1("user does not belong to game!")
		}

		if err != nil {
			log.Err(err).Msg("cannot send ExplosionEvent")
		}
	}
	return nil
}

func (r OutgoingEventHandlerImpl) EndGame(endGameEvent dto.EndGameEvent) error {
	if gameData, ok := cache.GameCache.Cache[utils.MaskId(endGameEvent.GameId)]; ok {

		eventBytes, err := dto.MarshalEvent(endGameEvent, dto.EndGame)
		if err != nil {
			log.Error().Err(err).Msg("cannot marshal EndGameEvent")
			return err
		}

		err = gameData.Side1Socket.WriteMessage(websocket.TextMessage, eventBytes)
		if err != nil {
			log.Err(err).Msg("cannot send EndGameEvent")
		}

		err = gameData.Side2Socket.WriteMessage(websocket.TextMessage, eventBytes)
		if err != nil {
			log.Err(err).Msg("cannot send EndGameEvent")
		}
	}
	return nil
}
