package dto

import (
	"battleship/model"
	"battleship/utils"
	"time"
)

type CreateGameRequest struct {
	UserId string `json:"user_id,omitempty"`
}

func (r *CreateGameRequest) GetUnMaskedUserId() string {
	return utils.MaskId(r.UserId)
}

type CreateGameResponse struct {
	BaseResponse
	GameId string `json:"game_id,omitempty"`
}

///////////////

type JoinGameRequest struct {
	GameId string `json:"game_id"`
	UserId string `json:"user_id"`
}

func (r *JoinGameRequest) GetUnMaskedUserId() string {
	return utils.MaskId(r.UserId)
}

func (r *JoinGameRequest) GetUnMaskedGameId() string {
	return utils.MaskId(r.GameId)
}

///////////////
type SubmitShipsLocationsRequest struct {
	UserId       string `json:"user_id"`
	GameId       string `json:"game_id"`
	ShipsIndexes []int  `json:"ships_indexes"`
}

func (r *SubmitShipsLocationsRequest) GetUnMaskedGameId() string {
	return utils.MaskId(r.GameId)
}

func (r *SubmitShipsLocationsRequest) GetUnMaskedUserId() string {
	return utils.MaskId(r.UserId)
}

//////////////

type MoveShipRequest struct {
}

//////////////

type GetGameResponse struct {
	BaseResponse
	Game *GameDto `json:"game"`
}

type GameDto struct {
	Id             string           `json:"id,omitempty"`
	State          *GameState       `json:"state,omitempty"`
	Status         model.GameStatus `json:"status,omitempty"`
	Side1UserId    *string          `json:"side_1_user_id,omitempty"`
	Side2UserId    *string          `json:"side_2_user_id,omitempty"`
	Turn           int              `json:"turn,omitempty"`
	LastMoveTime   *time.Time       `json:"last_move_time,omitempty"`
	MoveTimeoutSec int              `json:"move_timeout_sec,omitempty"`
	CreateDate     time.Time        `json:"create_date,omitempty"`
	Winner         *string          `json:"winner,omitempty"`
}

func (gameDto *GameDto) FromGame(game *model.Game) {
	gameDto.Id = utils.MaskId(game.Id.Hex())
	gameDto.Status = game.Status
	gameDto.LastMoveTime = game.LastMoveTime
	gameDto.Turn = game.Turn
	gameDto.MoveTimeoutSec = game.MoveTimeoutSec
	gameDto.CreateDate = game.CreateDate
	if game.Winner != nil {
		hex := game.Winner.Hex()
		gameDto.Winner = &hex
	}
	if game.Side1 != nil {
		hex := game.Side1.Hex()
		gameDto.Side1UserId = &hex
	}
	if game.Side2 != nil {
		hex := game.Side2.Hex()
		gameDto.Side2UserId = &hex
	}
	gameDto.State = &GameState{
		Side1:      game.State.Side1,
		Side2:      game.State.Side2,
		Side1Ships: game.State.Side1Ships,
		Side2Ships: game.State.Side2Ships,
	}
}

////////////////

type GameState struct {
	Side1      map[int]bool `json:"side_1,omitempty"`       //map index -> is hidden
	Side1Ships map[int]bool `json:"side_1_ships,omitempty"` //map index -> is ship exist
	Side2      map[int]bool `json:"side_2,omitempty"`       //map index -> is hidden
	Side2Ships map[int]bool `json:"side_2_ships,omitempty"` //map index -> is ship exist
}
