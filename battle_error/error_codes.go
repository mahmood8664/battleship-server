package battle_error

const (
	ShipInvalidMove ErrorCode = iota
	InvalidGameStatus
)

type ErrorCode int
