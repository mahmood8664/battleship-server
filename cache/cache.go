package cache

import (
	"battleship/dto"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"sync"
)

var GameCache = struct {
	Mux   sync.Locker
	Cache map[string]GameData
}{
	Cache: make(map[string]GameData),
	Mux:   new(sync.Mutex),
}

type GameData struct {
	Side1UserId string
	Side2UserId string
	Side1Socket *websocket.Conn
	Side2Socket *websocket.Conn
}

func (r *GameData) getUser(userId string) (string, error) {
	if r.Side1UserId == userId {
		return r.Side1UserId, nil
	} else if r.Side2UserId == userId {
		return r.Side2UserId, nil
	}
	log.Error().Str("userId", userId).Msg("cannot find userId in GameData")
	return "", dto.NotFoundError1(fmt.Sprintf("cannot find userId %s in GameData", userId))
}
