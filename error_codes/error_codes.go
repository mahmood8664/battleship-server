package error_codes

import "errors"

const (
	ShipInvalidMove ErrorCode = iota
	InvalidGameStatus
	InvalidShipIndexValue
	GameIsFinished
)

type ErrorCode int

var (
	ShipInvalidMoveRevealedLocation = errors.New("cannot move ship to revealed location")
	ShipInvalidMoveAlreadyDestroyed = errors.New("cannot move ship that is already destroyed")
	NotUserTurn                     = errors.New("its not user turn")
)
