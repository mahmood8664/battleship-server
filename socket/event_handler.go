package socket

import (
	"battleship/battle_error"
	"battleship/di"
	"battleship/dto"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var (
	gameService = di.CreateGameService()
)

func handleEvents(event dto.Event, c *websocket.Conn) error {

	switch event.Type {
	case dto.Connect:

		request := new(dto.SocketConnect)
		err := json.Unmarshal([]byte(event.Payload), request)
		if err != nil {
			return battle_error.ParseError(err)
		}

		gameDto, err := gameService.GetGame(request.GameId)
		if err != nil {
			return err
		}

		if _, ok := GameCache.cache[request.GameId]; ok == false {
			GameCache.mux.Lock()
			if _, ok := GameCache.cache[request.GameId]; ok == false {
				GameCache.cache[request.GameId] = GameData{}
			}
			GameCache.mux.Unlock()

			gameData := GameCache.cache[request.GameId]

			if gameDto.Side1UserId != nil && request.UserId == *gameDto.Side1UserId {
				gameData.Side1Socket = c
				gameData.Side1UserId = request.UserId
			} else if gameDto.Side2UserId != nil && request.UserId == *gameDto.Side2UserId {
				gameData.Side2Socket = c
				gameData.Side2UserId = request.UserId
			} else {
				log.Error().Str("user_id", request.UserId).
					Str("game_id", request.GameId).
					Msg("user does not belong to game")
				return battle_error.BadRequest1("user does not belong to game")
			}

		} else {
			gameData := GameCache.cache[request.GameId]

			if gameDto.Side1UserId != nil && request.UserId == *gameDto.Side1UserId {
				_ = gameData.Side1Socket.Close()
				gameData.Side1Socket = c
				gameData.Side1UserId = request.UserId
			} else if gameDto.Side2UserId != nil && request.UserId == *gameDto.Side2UserId {
				_ = gameData.Side2Socket.Close()
				gameData.Side2Socket = c
				gameData.Side2UserId = request.UserId
			} else {
				log.Error().Str("user_id", request.UserId).
					Str("game_id", request.GameId).
					Msg("user does not belong to game")
				return battle_error.BadRequest1("user does not belong to game")
			}
		}
	}
	return nil
}
