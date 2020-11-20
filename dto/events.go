package dto

import "battleship/utils"

const (
	Connect SocketEventType = iota + 1
)

type SocketEventType int

type Event struct {
	Type    SocketEventType `json:"event_type,omitempty"`
	Payload string          `json:"payload,omitempty"`
}

type SocketConnect struct {
	GameId string `json:"game_id"`
	UserId string `json:"user_id"`
}

func (r *SocketConnect) GetUnMaskedGameId() string {
	return utils.MaskId(r.GameId)
}

func (r *SocketConnect) GetUnMaskedUserId() string {
	return utils.MaskId(r.GameId)
}
