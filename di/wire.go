//+build wireinject

package di

import (
	"battleship/controllers"
	"battleship/db/dao"
	"battleship/service"
	"github.com/google/wire"
)

func CreateGameController() controllers.GameController {
	panic(wire.Build(
		controllers.NewGameControllerImpl,
		wire.Bind(new(controllers.GameController), new(*controllers.GameControllerImpl)),
		CreateGameService,
	))
}

func CreateUserController() controllers.UserController {
	panic(wire.Build(
		controllers.NewUserController,
		wire.Bind(new(controllers.UserController), new(*controllers.UserControllerImpl)),
		CreateUserService,
	))
}

func CreateGameService() service.GameService {
	panic(wire.Build(
		service.NewGameServiceImpl,
		wire.Bind(new(service.GameService), new(*service.GameServiceImpl)),
		CreateGameDao,
		CreateUserDao,
	))
}

func CreateUserService() service.UserService {
	panic(wire.Build(
		service.NewUserServiceImpl,
		wire.Bind(new(service.UserService), new(*service.UserServiceImpl)),
		CreateUserDao,
	))
}

func CreateGameDao() dao.GameDao {
	panic(wire.Build(
		dao.NewGameDaoImpl,
		wire.Bind(new(dao.GameDao), new(*dao.GameDaoImpl)),
	))
}

func CreateUserDao() dao.UserDao {
	panic(wire.Build(
		dao.NewUserDaoImpl,
		wire.Bind(new(dao.UserDao), new(*dao.UserDaoImpl)),
	))
}
