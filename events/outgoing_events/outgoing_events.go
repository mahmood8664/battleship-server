package outgoing_events

import (
	"battleship/dto"
	"github.com/rs/zerolog/log"
)

type OutgoingEventHandler interface {
	SendEvent(event dto.Event) error
}

type OutgoingEventHandlerImpl struct {
}

func NewOutgoingEventHandlerImpl() *OutgoingEventHandlerImpl {
	return &OutgoingEventHandlerImpl{}
}

func (r *OutgoingEventHandlerImpl) SendEvent(event dto.Event) error {
	switch event.Type {
	case dto.Connect:
		log.Error().Msg("Connect event is not allowed here")
	}
	return nil
}
