package dto

type UserDto struct {
	Id     string  `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
	Mobile *string `json:"mobile,omitempty"`
}

type CreateUserRequest struct {
	Name   *string `json:"name,omitempty"`
	Mobile *string `json:"mobile,omitempty"`
}

type CreateUserResponse struct {
	Id string `json:"id,omitempty"`
}
