package controllers

import (
	"battleship/battle_error"
	"battleship/dto"
	"battleship/service"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

type UserController interface {
	CreateUser() func(ctx echo.Context) error
	GetUser() func(ctx echo.Context) error
}

type UserControllerImpl struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserControllerImpl {
	return &UserControllerImpl{
		userService: userService,
	}
}

// Create user
// @Summary Create user
// @Description create a new user
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "Create User Request"
// @Success 200 {object} dto.CreateUserResponse "User successfully created"
// @Router /api/v1/user [post]
func (r *UserControllerImpl) CreateUser() func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		request := new(dto.CreateUserRequest)
		if err := ctx.Bind(request); err != nil {
			log.Warn().Err(err).Msg("Bad request")
			return battle_error.BadRequest1(err.Error())
		}

		res, err := r.userService.CreateUser(*request)
		if err != nil {
			log.Info().Err(err).Msg("cannot create user")
			return err
		}

		return ctx.JSON(http.StatusOK, res)
	}
}

// Get user
// @Summary Get user
// @Description Get user info
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path string true " "
// @Success 200 {object} dto.UserDto "Game successfully created"
// @Router /api/v1/user/{user_id} [get]
func (r *UserControllerImpl) GetUser() func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		userId := ctx.Param("user_id")
		if userId == "" {
			log.Warn().Msg("Bad request")
			return battle_error.BadRequest1("userId must has value")
		}

		user, err := r.userService.GetUser(userId)
		if err != nil {
			log.Info().Str("userId", userId).Err(err).Msg("cannot get user")
			return err
		}
		return ctx.JSON(http.StatusOK, user)
	}
}
