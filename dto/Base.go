package dto

type BaseRequest struct {
}

type UserGame interface {
	GetUserId() string
	GetGameId() string
}

type UserGameRequest struct {
	BaseRequest
	GameId string `json:"game_id"`
	UserId string `json:"user_id"`
}

func (r UserGameRequest) GetUserId() string {
	return r.UserId
}

func (r UserGameRequest) GetGameId() string {
	return r.GameId
}

type BaseResponse struct {
	Ok    bool         `json:"ok"`
	Error *BattleError `json:"error,omitempty"`
}
