package dto

import (
	"battleship/utils"
	"encoding/json"
	"github.com/rs/zerolog/log"
)

const (
	Connect    SocketEventType = "connect"
	GameStart                  = "game_start"
	ChangeTurn                 = "change_turn"
	ShipMoved                  = "ship_moved"
	Reveal                     = "reveal"
	Explosion                  = "explosion"
	EndGame                    = "end_game"
)

type SocketEventType string

type Event struct {
	Type    SocketEventType `json:"event_type,omitempty"`
	Payload string          `json:"payload,omitempty"`
}

func MarshalEvent(payload interface{}, eventType SocketEventType) ([]byte, error) {
	marshal, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Str("event_type", string(eventType)).Msg("cannot marshal payload")
		return nil, err
	}
	gameStartEventBytes := Event{
		Type:    eventType,
		Payload: string(marshal),
	}
	eventBytes, err := json.Marshal(gameStartEventBytes)
	if err != nil {
		log.Error().Err(err).Str("event_type", string(eventType)).Msg("cannot marshal Event")
		return nil, err
	}
	return eventBytes, nil
}

/////////////

type UserConnectEvent struct {
	UserGameRequest
}

func (r *UserConnectEvent) ValidateAndUnmask() error {
	if r.UserId == "" {
		return BadRequest1("user id is not correct")
	}
	if r.GameId == "" {
		return BadRequest1("game id is not correct")
	}
	r.UserId = utils.MaskId(r.UserId)
	r.GameId = utils.MaskId(r.GameId)
	return nil
}

//////////////

type GameConnect struct {
	GameId string `json:"game_id"`
	UserId string `json:"user_id"`
}

//////////////

type GameStartEvent struct {
	Game GameDto `json:"game"`
}

//////////////

type GameChangeTurnEvent struct {
	GameId string `json:"game_id"`
	UserId string `json:"user_id"`
}

//////////////

type ShipMovedEvent struct {
	GameId       string `json:"game_id"`
	UserId       string `json:"user_id"`
	OldShipIndex int    `json:"old_ship_index"`
}

//////////

type RevealEvent struct {
	GameId        string `json:"game_id"`
	UserId        string `json:"user_id"`
	Slots         []int  `json:"slots"`
	RevealedShips []int  `json:"revealed_ships"`
}

////////////
type ExplosionEvent struct {
	GameId string `json:"game_id"`
	UserId string `json:"user_id"`
	Index  int    `json:"index"`
}

////////////
type EndGameEvent struct {
	GameId       string `json:"game_id"`
	WinnerUserId string `json:"winner_user_id"`
}
