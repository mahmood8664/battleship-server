package dto

type BaseRequest struct {
}

type BaseResponse struct {
	Ok    bool         `json:"ok,omitempty"`
	Error *BattleError `json:"error,omitempty"`
}
