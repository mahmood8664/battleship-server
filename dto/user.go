package dto

type UserDto struct {
	BaseResponse
	Id     string  `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
	Mobile *string `json:"mobile,omitempty"`
}

type CreateUserRequest struct {
	Name   *string `json:"name,omitempty"`
	Mobile *string `json:"mobile,omitempty"`
}

type CreateUserResponse struct {
	BaseResponse
	Id string `json:"id,omitempty"`
}
