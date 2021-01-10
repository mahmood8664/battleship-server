package dto

import (
	"battleship/model"
	"battleship/utils"
	"github.com/rs/zerolog/log"
	"time"
)

type CreateGameRequest struct {
	UserId      string `json:"user_id,omitempty"`
	MoveTimeout int    `json:"move_timeout"`
}

func (r *CreateGameRequest) ValidateAndUnmask() error {
	if r.UserId == "" {
		return BadRequest1("user id is not correct")
	}
	if r.MoveTimeout < 5 && r.MoveTimeout > 30 {
		return BadRequest1("move timeout is between 5 and 30")
	}
	r.UserId = utils.MaskId(r.UserId)
	return nil
}

///////////////

type JoinGameRequest struct {
	UserGameRequest
}

func (r *JoinGameRequest) ValidateAndUnmask() error {
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

///////////////
type SubmitShipsLocationsRequest struct {
	UserGameRequest
	ShipsIndexes []int `json:"ships_indexes"`
}

func (r *SubmitShipsLocationsRequest) ValidateAndUnmask() error {
	if r.UserId == "" {
		return BadRequest1("user id is not correct")
	}
	if r.GameId == "" {
		return BadRequest1("game id is not correct")
	}
	if len(r.ShipsIndexes) != 10 {
		log.Error().Msg("ship index size must be 10")
		return BadRequest1("ship index size must be 10")
	}
	r.UserId = utils.MaskId(r.UserId)
	r.GameId = utils.MaskId(r.GameId)
	return nil
}

type SubmitShipsLocationsResponse struct {
	BaseResponse
	GameStatue model.GameStatus `json:"game_status"`
	Turn       int              `json:"turn"`
}

//////////////

type MoveShipRequest struct {
	UserGameRequest
	OldShipIndex int `json:"old_ship_index"`
	NewShipIndex int `json:"new_ship_index"`
}

func (r *MoveShipRequest) ValidateAndUnmask() error {
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

type MoveShipResponse struct {
	BaseResponse
}

//////////////

type ChangeTurnRequest struct {
	UserGameRequest
}

func (r *ChangeTurnRequest) ValidateAndUnmask() error {
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

type ChangeTurnResponse struct {
	BaseResponse
}

////////////

type RevealEnemyFieldsRequest struct {
	UserGameRequest
	Index int `json:"index"`
}

func (r *RevealEnemyFieldsRequest) ValidateAndUnmask() error {
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

type RevealEnemyFieldsResponse struct {
	BaseResponse
	RevealedShipIndexes []int `json:"revealed_ship_indexes"`
}

////////////

type ExplodeRequest struct {
	UserGameRequest
	Index int `json:"index"`
}

func (r *ExplodeRequest) ValidateAndUnmask() error {
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

type ExplodeResponse struct {
	BaseResponse
	HasShip bool `json:"has_ship"`
}

////////////

type GetGameRequest struct {
	UserGameRequest
}

type GetGameResponse struct {
	BaseResponse
	Game *GameDto `json:"game"`
}

type GameDto struct {
	Id              string           `json:"id,omitempty"`
	State           *GameState       `json:"state,omitempty"`
	Status          model.GameStatus `json:"status,omitempty"`
	UserId          string           `json:"user_id,omitempty"`
	YourTurn        bool             `json:"your_turn"`
	OtherSideJoined bool             `json:"other_side_joined"`
	MoveTimeoutSec  int              `json:"move_timeout_sec,omitempty"`
	CreateDate      time.Time        `json:"create_date,omitempty"`
	WinnerUser      *string          `json:"winner_user,omitempty"`
}

type GameState struct {
	OwnGround          map[int]bool `json:"own_ground,omitempty"`
	OwnShips           map[int]bool `json:"own_ships,omitempty"`
	EnemyGround        map[int]bool `json:"enemy_ground,omitempty"`
	EnemyRevealedShips map[int]bool `json:"enemy_revealed_ships,omitempty"`
}

func (r *GameDto) FromGame(game model.Game, requesterUserId string) {
	r.Id = utils.MaskId(game.Id.Hex())
	r.Status = game.Status
	r.MoveTimeoutSec = game.MoveTimeoutSec
	r.CreateDate = game.CreateDate
	if game.WinnerUser != nil {
		winnerId := utils.MaskId(game.WinnerUser.Hex())
		r.WinnerUser = &winnerId
	}

	if game.Side1User != nil && requesterUserId == game.Side1User.Hex() {
		r.UserId = utils.MaskId(game.Side1User.Hex())
		if game.Turn == 1 {
			r.YourTurn = true
		} else {
			r.YourTurn = false
		}

		r.State = &GameState{
			OwnGround:          game.State.Side1Ground,
			OwnShips:           game.State.Side1Ships,
			EnemyGround:        game.State.Side2Ground,
			EnemyRevealedShips: game.State.Side2RevealedShips,
		}

		if game.Side2User != nil {
			r.OtherSideJoined = true
		}
	} else if game.Side2User != nil && requesterUserId == game.Side2User.Hex() {
		r.UserId = utils.MaskId(game.Side2User.Hex())
		if game.Turn == 2 {
			r.YourTurn = true
		} else {
			r.YourTurn = false
		}

		r.State = &GameState{
			OwnGround:          game.State.Side2Ground,
			OwnShips:           game.State.Side2Ships,
			EnemyGround:        game.State.Side1Ground,
			EnemyRevealedShips: game.State.Side1RevealedShips,
		}

		if game.Side1User != nil {
			r.OtherSideJoined = true
		}
	} else {
		log.Error().Str("user_id", requesterUserId).Str("game_id", game.Id.Hex()).Msg("user does not belong to this game")
	}
}
