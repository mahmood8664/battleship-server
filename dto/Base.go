package dto

type BaseRequest struct {
	ReqId string `json:"req_id,omitempty"`
}

type BaseResponse struct {
	ReqId        string `json:"req_id,omitempty"`
	Successful   bool   `json:"successful,omitempty"`
	ErrorCode    int    `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}
