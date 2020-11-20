package incoming_events

import (
	"battleship/dto"
	"battleship/service"
	"github.com/rs/zerolog/log"
)

type IncomingEventHandler interface {
	HandleEvent(event dto.Event) error
}

type IncomingEventHandlerImpl struct {
	gameService service.GameService
}

func NewIncomingEventHandlerImpl(gameService service.GameService) *IncomingEventHandlerImpl {
	return &IncomingEventHandlerImpl{
		gameService: gameService,
	}
}

func (r *IncomingEventHandlerImpl) HandleEvent(event dto.Event) error {
	switch event.Type {
	case dto.Connect:
		log.Error().Msg("Connect event is not allowed here")
	}
	return nil
}
