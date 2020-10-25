package battle_error

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

type BattleError struct {
	Cause     error             `json:"cause,omitempty"`
	HttpCode  int               `json:"httpCode,omitempty"`
	ErrorCode ErrorCode         `json:"errorCode,omitempty"`
	Message   string            `json:"message,omitempty"`
	Data      map[string]string `json:"data,omitempty"`
}

func (e *BattleError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("BattleError with message '%s' and http code '%d' and server code %d and cause '%s'",
			e.Message, e.HttpCode, e.ErrorCode, e.Cause.Error())
	} else {
		return fmt.Sprintf("BattleError with message '%s' and http code '%d' and server code %d",
			e.Message, e.HttpCode, e.ErrorCode)
	}
}

//////////////////////////////////////////////////////////
// Define your custom error HERE:

func ParseError(err error) error {
	if err == nil {
		return nil
	}

	var e = &BattleError{}
	if errors.As(err, &e) {
		return &BattleError{
			Cause:    err,
			HttpCode: err.(*BattleError).HttpCode,
			Message:  err.(*BattleError).Message,
			Data:     err.(*BattleError).Data,
		}
	} else {
		return parseErrorMessage(err)
	}
}

func parseErrorMessage(err error) error {

	switch {
	case err.Error() == "mongo: no documents in result":
		return &BattleError{
			Cause:    err,
			HttpCode: http.StatusNotFound,
			Message:  "not found",
			Data:     nil,
		}
	case strings.HasPrefix(err.Error(), "encoding/hex:"):
		return BadRequest1("Invalid id")
	default:
		return &BattleError{
			Cause:    err,
			HttpCode: http.StatusInternalServerError,
			Message:  "cannot parse error message: " + err.Error(),
			Data:     nil,
		}
	}
}

func EntityNotFound(message string, entityName string, entityValue string) error {
	return &BattleError{
		HttpCode: http.StatusNotFound,
		Message:  message,
		Data:     map[string]string{"entityName": entityName, "entityValue": entityValue},
	}
}

func ValidationError(message string) error {
	return &BattleError{
		HttpCode: http.StatusBadRequest,
		Message:  message,
	}
}

func BadRequest0() error {
	return &BattleError{
		HttpCode: http.StatusBadRequest,
	}
}

func BadRequest1(message string) error {
	return &BattleError{
		HttpCode: http.StatusBadRequest,
		Message:  message,
	}
}

func BadRequest2(message string, code ErrorCode) error {
	return &BattleError{
		HttpCode:  http.StatusBadRequest,
		ErrorCode: code,
		Message:   message,
	}
}

func Unauthorized(message string) error {
	return &BattleError{
		HttpCode: http.StatusUnauthorized,
		Message:  message,
	}
}

func Forbidden(message string, errorCode ErrorCode) error {
	return &BattleError{
		HttpCode:  http.StatusForbidden,
		ErrorCode: errorCode,
		Message:   message,
	}
}

func NotFoundError0() error {
	return &BattleError{
		HttpCode: http.StatusNotFound,
	}
}

func NotFoundError1(message string) error {
	return &BattleError{
		HttpCode: http.StatusNotFound,
		Message:  message,
	}
}

func NotFoundError2(message string, errorCode ErrorCode) error {
	return &BattleError{
		HttpCode:  http.StatusNotFound,
		ErrorCode: errorCode,
		Message:   message,
	}
}

/////////////////////////////////////////////////////////
func CustomHTTPErrorHandler(err error, c echo.Context) {
	if e, ok := err.(*BattleError); ok {
		handleError(c, e)
	} else if e, ok := err.(*echo.HTTPError); ok {
		handleHttpError(c, e)
	} else {
		err := ParseError(err)
		handleError(c, err.(*BattleError))
	}
}

func handleError(c echo.Context, e *BattleError) {
	//log.Warn().Msg(e.Message)
	if err2 := c.JSON(e.HttpCode, e); err2 != nil {
		log.Error().Msg("error in converting BattleError to json: " + err2.Error())
	}
}

func handleHttpError(c echo.Context, e *echo.HTTPError) {
	//log.Warn().Msg(cast.ToString(e.Message))
	if err2 := c.JSON(e.Code, e); err2 != nil {
		log.Error().Msg("error in converting HttpError to json: " + err2.Error())
	}
}
