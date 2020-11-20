package model

import (
	"battleship/error_codes"
	"battleship/utils"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type GameStatus string

const (
	Init     GameStatus = "init"
	Start    GameStatus = "start"
	Finished GameStatus = "finished"
)

type Game struct {
	Id             primitive.ObjectID  `bson:"_id,omitempty"`
	State          GameState           `bson:"state,omitempty"`
	Status         GameStatus          `bson:"status,omitempty"`
	Side1          *primitive.ObjectID `bson:"side_1,omitempty"`
	Side2          *primitive.ObjectID `bson:"side_2,omitempty"`
	Turn           int                 `bson:"turn,omitempty"`
	LastMoveTime   *time.Time          `bson:"last_move_time,omitempty"`
	MoveTimeoutSec int                 `bson:"move_timeout_sec,omitempty"`
	CreateDate     time.Time           `bson:"create_date,omitempty"`
	Winner         *primitive.ObjectID `bson:"winner,omitempty"`
}

func (g *Game) GetMaskedGameId() string {
	return utils.MaskId(g.Id.Hex())
}

func (g *Game) moveShipSide1(shipIndex int, locationIndex int) error {
	if exist, ok := g.State.Side1Ships[shipIndex]; ok && exist {
		if val, ok2 := g.State.Side1[locationIndex]; ok2 && val {
			delete(g.State.Side1Ships, shipIndex)
			g.State.Side1Ships[locationIndex] = true
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

func (g *Game) moveShipSide2(shipIndex int, locationIndex int) error {
	if exist, ok := g.State.Side2Ships[shipIndex]; ok && exist {
		if val, ok2 := g.State.Side2[locationIndex]; ok2 && val {
			delete(g.State.Side2Ships, shipIndex)
			g.State.Side2Ships[locationIndex] = true
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

func (g Game) RevealSlot1(index int) (notEmptySlots []int) {
	neighborIndexes := findNeighborIndexes(index)
	for i := range neighborIndexes {
		g.State.Side1Ships[i] = false
		if g.State.Side1Ships[i] {
			notEmptySlots = append(notEmptySlots, i)
		}
	}
	return notEmptySlots
}

func (g Game) RevealSlot2(index int) (notEmptySlots []int) {
	neighborIndexes := findNeighborIndexes(index)
	for i := range neighborIndexes {
		g.State.Side2Ships[i] = false
		if g.State.Side2Ships[i] {
			notEmptySlots = append(notEmptySlots, i)
		}
	}
	return notEmptySlots
}

func (g Game) ExplodeSide1(index int) {
	if g.State.Side1Ships[index] {
		delete(g.State.Side1Ships, index)
		if len(g.State.Side1Ships) == 0 {
			g.Status = Finished
			g.Winner = g.Side2
		}
	}
}

func (g Game) ExplodeSide2(index int) {
	if g.State.Side2Ships[index] {
		delete(g.State.Side2Ships, index)
		if len(g.State.Side2Ships) == 0 {
			g.Status = Finished
			g.Winner = g.Side1
		}
	}
}

func findNeighborIndexes(index int) [4]int {
	index1 := index + 1
	index2 := index + 20
	index3 := index + 21
	if index > 0 && (index+1)%20 == 0 {

		index1 = index - 1
		index2 = index + 19
		index3 = index + 20
	}
	if index+1 > 9*20 {

		index1 = index + 1
		index2 = index - 19
		index3 = index - 20
	}
	if index == 199 {

		index1 = 198
		index2 = 199 - 20
		index3 = 199 - 21
	}
	return [4]int{index, index1, index2, index3}
}
