package model

import (
	"battleship/error_codes"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type GameStatus string

const (
	Init     GameStatus = "init"
	Joined   GameStatus = "joined"
	Start    GameStatus = "start"
	Finished GameStatus = "finished"
)

type Game struct {
	Id             primitive.ObjectID  `bson:"_id,omitempty"`
	State          GameState           `bson:"state"`
	LastMoveTime   time.Time           `bson:"last_move_time"`
	Status         GameStatus          `bson:"status"`
	Side1User      *primitive.ObjectID `bson:"side_1_user"`
	Side2User      *primitive.ObjectID `bson:"side_2_user"`
	Turn           int                 `bson:"turn"`
	MoveTimeoutSec int                 `bson:"move_timeout_sec"`
	CreateDate     time.Time           `bson:"create_date"`
	WinnerUser     *primitive.ObjectID `bson:"winner_user"`
}

func (g *Game) MoveShipSide1(from int, to int) error {
	if exist, ok := g.State.Side1Ships[from]; ok && exist {
		if val, ok2 := g.State.Side1Ground[to]; ok2 && val {
			delete(g.State.Side1Ships, from)
			g.State.Side1Ships[to] = true
			delete(g.State.Side1RevealedShips, from)
		} else {
			log.Debug().Str("gameId", g.Id.Hex()).Msg("cannot move ship to revealed location")
			return error_codes.ShipInvalidMoveRevealedLocation
		}
	} else {
		log.Debug().Str("gameId", g.Id.Hex()).Msg("cannot move ship that is already destroyed")
		return error_codes.ShipInvalidMoveAlreadyDestroyed
	}
	return nil
}

func (g *Game) MoveShipSide2(from int, to int) error {
	if exist, ok := g.State.Side2Ships[from]; ok && exist {
		if val, ok2 := g.State.Side2Ground[to]; ok2 && val {
			delete(g.State.Side2Ships, from)
			g.State.Side2Ships[to] = true
			delete(g.State.Side2RevealedShips, from)
		} else {
			log.Debug().Str("gameId", g.Id.Hex()).Msg("cannot move ship to revealed location")
			return error_codes.ShipInvalidMoveRevealedLocation
		}
	} else {
		log.Debug().Str("gameId", g.Id.Hex()).Msg("cannot move ship that is already destroyed")
		return error_codes.ShipInvalidMoveAlreadyDestroyed
	}
	return nil
}

func (g *Game) RevealSlotSide1(index int) (notEmptySlots []int) {
	neighborIndexes := FindNeighborIndexes(index)
	for _, i := range neighborIndexes {
		if g.State.Side1Ships[i] {
			notEmptySlots = append(notEmptySlots, i)
			g.State.Side1RevealedShips[i] = true
		}
		g.State.Side1Ground[i] = false
	}
	return notEmptySlots
}

func (g *Game) RevealSlotSide2(index int) (notEmptySlots []int) {
	neighborIndexes := FindNeighborIndexes(index)
	for _, i := range neighborIndexes {
		if g.State.Side2Ships[i] {
			notEmptySlots = append(notEmptySlots, i)
			g.State.Side2RevealedShips[i] = true
		}
		g.State.Side2Ground[i] = false
	}
	return notEmptySlots
}

func (g *Game) ExplodeSide1(index int) bool {
	if g.State.Side1Ships[index] {
		g.State.Side1Ships[index] = false
		if allShipsDestroyed(g.State.Side1Ships) {
			g.Status = Finished
			g.WinnerUser = g.Side2User
		}
		return true
	}
	return false
}

func (g *Game) ExplodeSide2(index int) bool {
	if g.State.Side2Ships[index] {
		g.State.Side2Ships[index] = false
		if allShipsDestroyed(g.State.Side2Ships) {
			g.Status = Finished
			g.WinnerUser = g.Side1User
		}
		return true
	}
	return false
}

func allShipsDestroyed(ships map[int]bool) bool {
	for _, b := range ships {
		if b {
			return false
		}
	}
	return true
}

func FindNeighborIndexes(index int) []int {
	index1 := index + 1
	index2 := index + 10
	index3 := index + 11
	if index > 0 && (index+1)%10 == 0 {

		index1 = index - 1
		index2 = index + 9
		index3 = index + 10
	}
	if index+1 > 9*10 {

		index1 = index + 1
		index2 = index - 9
		index3 = index - 10
	}
	if index == 99 {

		index1 = 98
		index2 = 99 - 10
		index3 = 99 - 11
	}
	return []int{index, index1, index2, index3}
}
