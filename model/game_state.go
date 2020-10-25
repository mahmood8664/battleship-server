package model

type GameState struct {
	Side1      map[int]bool `bson:"side_1,omitempty"`       //map index -> is hidden
	Side1Ships map[int]bool `bson:"side_1_ships,omitempty"` //map index -> is ship exist
	Side2      map[int]bool `bson:"side_2,omitempty"`       //map index -> is hidden
	Side2Ships map[int]bool `bson:"side_2_ships,omitempty"` //map index -> is ship exist
}

func CreateGameState(side1ShipsIndexes []int, side2ShipsIndexes []int) *GameState {

	side1Ships := make(map[int]bool)
	side2Ships := make(map[int]bool)
	side1 := make(map[int]bool)
	side2 := make(map[int]bool)

	for i := range side1ShipsIndexes {
		side1Ships[i] = true
	}

	for i := range side2ShipsIndexes {
		side2Ships[i] = true
	}

	for i := 0; i < 200; i++ {
		side1[i] = true
		side2[i] = true
	}

	return &GameState{
		Side1:      side1,
		Side1Ships: side1Ships,
		Side2:      side2,
		Side2Ships: side2Ships,
	}
}
