package incoming_events

import (
	"battleship/db/dao"
	"battleship/dto"
)

type IncomingEventHandler interface {
	HandleEvent(event dto.Event) error
}

type IncomingEventHandlerImpl struct {
	gameDao dao.GameDao
}

func NewIncomingEventHandlerImpl(gameDao dao.GameDao) IncomingEventHandlerImpl {
	return IncomingEventHandlerImpl{
		gameDao: gameDao,
	}
}

func (r IncomingEventHandlerImpl) HandleEvent(event dto.Event) error {
	return nil
}
