//+build wireinject

package di

import (
	"battleship/controllers"
	"battleship/db/dao"
	"battleship/events/incoming_events"
	"battleship/events/outgoing_events"
	"battleship/service"
	"battleship/socket"
	"github.com/google/wire"
)

func CreateGameController() controllers.GameController {
	panic(wire.Build(
		controllers.NewGameControllerImpl,
		wire.Bind(new(controllers.GameController), new(controllers.GameControllerImpl)),
		CreateGameService,
	))
}

func CreateUserController() controllers.UserController {
	panic(wire.Build(
		controllers.NewUserController,
		wire.Bind(new(controllers.UserController), new(controllers.UserControllerImpl)),
		CreateUserService,
	))
}

//////////////

func CreateGameService() service.GameService {
	panic(wire.Build(
		service.NewGameServiceImpl,
		wire.Bind(new(service.GameService), new(service.GameServiceImpl)),
		CreateGameDao,
		CreateUserDao,
		CreateGameEventDao,
		CreateOutgoingEventHandler,
	))
}

func CreateUserService() service.UserService {
	panic(wire.Build(
		service.NewUserServiceImpl,
		wire.Bind(new(service.UserService), new(service.UserServiceImpl)),
		CreateUserDao,
	))
}

func CreateIncomingEventHandler() incoming_events.IncomingEventHandler {
	panic(wire.Build(
		incoming_events.NewIncomingEventHandlerImpl,
		wire.Bind(new(incoming_events.IncomingEventHandler), new(incoming_events.IncomingEventHandlerImpl)),
		CreateGameDao,
	))
}

func CreateOutgoingEventHandler() outgoing_events.OutgoingEventHandler {
	panic(wire.Build(
		outgoing_events.NewOutgoingEventHandlerImpl,
		wire.Bind(new(outgoing_events.OutgoingEventHandler), new(outgoing_events.OutgoingEventHandlerImpl)),
	))
}

func CreateSocketHandler() socket.SocketHandler {
	panic(wire.Build(
		socket.NewSocketHandlerImpl,
		wire.Bind(new(socket.SocketHandler), new(socket.SocketHandlerImpl)),
		CreateIncomingEventHandler,
		CreateGameService,
	))
}

////////////

func CreateGameDao() dao.GameDao {
	panic(wire.Build(
		dao.NewGameDaoImpl,
		wire.Bind(new(dao.GameDao), new(dao.GameDaoImpl)),
	))
}

func CreateUserDao() dao.UserDao {
	panic(wire.Build(
		dao.NewUserDaoImpl,
		wire.Bind(new(dao.UserDao), new(dao.UserDaoImpl)),
	))
}

func CreateGameEventDao() dao.GameEventDao {
	panic(wire.Build(
		dao.NewEventGameDaoImpl,
		wire.Bind(new(dao.GameEventDao), new(dao.GameEventDaoImpl)),
	))
}
