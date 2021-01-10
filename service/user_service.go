package service

import (
	"battleship/db/dao"
	"battleship/dto"
	"battleship/model"
	"battleship/utils"
	"github.com/rs/zerolog/log"
)

type UserService interface {
	CreateUser(request dto.CreateUserRequest) (response dto.CreateUserResponse, err error)
	GetUser(id string) (user dto.UserDto, err error)
}

type UserServiceImpl struct {
	userDao dao.UserDao
}

func NewUserServiceImpl(userDao dao.UserDao) UserServiceImpl {
	return UserServiceImpl{
		userDao: userDao,
	}
}

func (r UserServiceImpl) CreateUser(request dto.CreateUserRequest) (response dto.CreateUserResponse, err error) {
	id, err := r.userDao.Insert(model.User{
		Name:   request.Name,
		Mobile: request.Mobile,
	})
	if err == nil {
		response.Id = utils.MaskId(id)
		response.Ok = true
	} else {
		log.Info().Err(err).Msg("cannot create user")
	}
	return response, err
}

func (r UserServiceImpl) GetUser(id string) (user dto.UserDto, err error) {
	u, err := r.userDao.GetOne(id)
	if err == nil {
		user.Mobile = u.Mobile
		user.Name = u.Name
		user.Id = utils.MaskId(u.Id.Hex())
		user.Ok = true
	} else {
		log.Info().Str("userId", id).Err(err).Msg("cannot get user")
	}
	return user, err
}
