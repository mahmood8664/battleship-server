package error_codes

import "errors"

const (
	ShipInvalidMove ErrorCode = iota
	InvalidGameStatus
	InvalidShipIndexValue
)

type ErrorCode int

var (
	ShipInvalidMoveRevealedLocation = errors.New("cannot move ship to revealed location")
	ShipInvalidMoveAlreadyDestroyed = errors.New("cannot move ship that is already destroyed")
)
