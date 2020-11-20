package controllers

import (
	"battleship/dto"
	"battleship/service"
	"battleship/utils"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

type UserController interface {
	CreateUser(ctx echo.Context) error
	GetUser(ctx echo.Context) error
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
func (r *UserControllerImpl) CreateUser(ctx echo.Context) error {
	request := new(dto.CreateUserRequest)
	if err := ctx.Bind(request); err != nil {
		log.Warn().Err(err).Msg("Bad request")
		return dto.BadRequest1(err.Error())
	}

	res, err := r.userService.CreateUser(*request)
	if err != nil {
		log.Info().Err(err).Msg("cannot create user")
		return err
	}
	return ctx.JSON(http.StatusOK, res)
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
func (r *UserControllerImpl) GetUser(ctx echo.Context) error {
	userId := ctx.Param("user_id")
	if userId == "" {
		log.Warn().Msg("Bad request")
		return dto.BadRequest1("userId must has value")
	}

	user, err := r.userService.GetUser(utils.MaskId(userId))
	if err != nil {
		log.Info().Str("userId", utils.MaskId(userId)).Err(err).Msg("cannot get user")
		return err
	}
	return ctx.JSON(http.StatusOK, user)
}
