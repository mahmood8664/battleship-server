package model

type GameState struct {
	Side1Ground        map[int]bool `bson:"side_1_ground"`         //map index -> is hidden
	Side1Ships         map[int]bool `bson:"side_1_ships"`          //map index -> is ship exist
	Side1RevealedShips map[int]bool `bson:"side_1_revealed_ships"` //map index -> is ship exploded
	Side2Ground        map[int]bool `bson:"side_2_ground"`         //map index -> is hidden
	Side2Ships         map[int]bool `bson:"side_2_ships"`          //map index -> is ship exist
	Side2RevealedShips map[int]bool `bson:"side_2_revealed_ships"` //map index -> is ship exploded
}
