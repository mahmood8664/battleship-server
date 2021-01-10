package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	JoinGame              GameEventType = "join_game"
	InitialShipsLocations               = "initial_ship_location"
	MoveShip                            = "move_ship"
	Explosion                           = "explosion"
	EmptyExplosion                      = "empty_explosion"
	ChangeTurn                          = "change_turn"
	Reveal                              = "reveal"
)

type GameEventType string

type GameEvent struct {
	Id                    primitive.ObjectID  `bson:"_id,omitempty"`
	Type                  GameEventType       `bson:"type"`
	InitialShipsLocations []int               `bson:"initial_ships_locations,omitempty"`
	MoveShipFrom          *int                `bson:"move_ship_from,omitempty"`
	MoveShipTo            *int                `bson:"move_ship_to,omitempty"`
	DiscoverEnemy         []int               `bson:"discover_enemy,omitempty"`
	DiscoverEnemyShips    []int               `bson:"discover_enemy_ships,omitempty"`
	Explosion             *int                `bson:"explosion,omitempty"`
	EmptyExplosion        *int                `bson:"empty_explosion,omitempty"`
	Time                  time.Time           `bson:"time,omitempty"`
	UserId                *primitive.ObjectID `bson:"user_id,omitempty"`
	GameId                primitive.ObjectID  `bson:"game_id,omitempty"`
}
