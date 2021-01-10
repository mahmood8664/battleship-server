package dto

import (
	"battleship/error_codes"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

type BattleError struct {
	ErrorCause    error                 `json:"-"`
	HttpErrorCode int                   `json:"-"`
	ErrorCode     error_codes.ErrorCode `json:"error_code,omitempty"`
	ErrorMessage  string                `json:"error_message,omitempty"`
	ErrorData     map[string]string     `json:"error_data,omitempty"`
}

func (e *BattleError) Error() string {
	if e.ErrorCause != nil {
		return fmt.Sprintf("BattleError with message '%s' and http code '%d' and server code %d and cause '%s'",
			e.ErrorMessage, e.HttpErrorCode, e.ErrorCode, e.ErrorCause.Error())
	} else {
		return fmt.Sprintf("BattleError with message '%s' and http code '%d' and server code %d",
			e.ErrorMessage, e.HttpErrorCode, e.ErrorCode)
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
			ErrorCause:    err,
			HttpErrorCode: err.(*BattleError).HttpErrorCode,
			ErrorMessage:  err.(*BattleError).ErrorMessage,
			ErrorData:     err.(*BattleError).ErrorData,
		}
	} else {
		return parseErrorMessage(err)
	}
}

func parseErrorMessage(err error) error {

	switch {
	case err.Error() == "mongo: no documents in result":
		return &BattleError{
			ErrorCause:    err,
			HttpErrorCode: http.StatusNotFound,
			ErrorMessage:  "not found",
			ErrorData:     nil,
		}
	case err.Error() == error_codes.ShipInvalidMoveRevealedLocation.Error():
		return BadRequest2(err.Error(), error_codes.ShipInvalidMove)
	case err.Error() == error_codes.ShipInvalidMoveAlreadyDestroyed.Error():
		return BadRequest2(err.Error(), error_codes.ShipInvalidMove)
	case err.Error() == error_codes.NotUserTurn.Error():
		return Forbidden2(err.Error(), error_codes.ShipInvalidMove)
	case strings.HasPrefix(err.Error(), "encoding/hex:"):
		return BadRequest1("Invalid id")
	default:
		return &BattleError{
			ErrorCause:    err,
			HttpErrorCode: http.StatusInternalServerError,
			ErrorMessage:  "cannot parse error message: " + err.Error(),
			ErrorData:     nil,
		}
	}
}

func EntityNotFound(message string, entityName string, entityValue string) error {
	return &BattleError{
		HttpErrorCode: http.StatusNotFound,
		ErrorMessage:  message,
		ErrorData:     map[string]string{"entityName": entityName, "entityValue": entityValue},
	}
}

func ValidationError(message string) error {
	return &BattleError{
		HttpErrorCode: http.StatusBadRequest,
		ErrorMessage:  message,
	}
}

func BadRequest0() error {
	return &BattleError{
		HttpErrorCode: http.StatusBadRequest,
		ErrorMessage:  "Bad request",
	}
}

func BadRequest1(message string) error {
	return &BattleError{
		HttpErrorCode: http.StatusBadRequest,
		ErrorMessage:  message,
	}
}

func BadRequest2(message string, code error_codes.ErrorCode) error {
	return &BattleError{
		HttpErrorCode: http.StatusBadRequest,
		ErrorCode:     code,
		ErrorMessage:  message,
	}
}

func Unauthorized(message string) error {
	return &BattleError{
		HttpErrorCode: http.StatusUnauthorized,
		ErrorMessage:  message,
	}
}

func Forbidden1(message string) error {
	return &BattleError{
		HttpErrorCode: http.StatusForbidden,
		ErrorMessage:  message,
	}
}

func Forbidden2(message string, code error_codes.ErrorCode) error {
	return &BattleError{
		HttpErrorCode: http.StatusForbidden,
		ErrorMessage:  message,
		ErrorCode:     code,
	}
}

func NotFoundError0() error {
	return &BattleError{
		HttpErrorCode: http.StatusNotFound,
		ErrorMessage:  "Not found",
	}
}

func NotFoundError1(message string) error {
	return &BattleError{
		HttpErrorCode: http.StatusNotFound,
		ErrorMessage:  message,
	}
}

func NotFoundError2(message string, errorCode error_codes.ErrorCode) error {
	return &BattleError{
		HttpErrorCode: http.StatusNotFound,
		ErrorCode:     errorCode,
		ErrorMessage:  message,
	}
}

func Duplicate0() error {
	return &BattleError{
		HttpErrorCode: http.StatusConflict,
		ErrorMessage:  "duplicate request",
	}
}

func Duplicate1(message string) error {
	return &BattleError{
		HttpErrorCode: http.StatusConflict,
		ErrorMessage:  message,
	}
}

func Duplicate11(code error_codes.ErrorCode) error {
	return &BattleError{
		HttpErrorCode: http.StatusConflict,
		ErrorCode:     code,
	}
}

func Duplicate2(message string, code error_codes.ErrorCode) error {
	return &BattleError{
		HttpErrorCode: http.StatusConflict,
		ErrorMessage:  message,
		ErrorCode:     code,
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
	response := BaseResponse{
		Ok:    false,
		Error: e,
	}
	if err2 := c.JSON(e.HttpErrorCode, response); err2 != nil {
		log.Error().Msg("error in converting BattleError to json: " + err2.Error())
	}
}

func handleHttpError(c echo.Context, e *echo.HTTPError) {
	//log.Warn().Msg(cast.ToString(e.Message))
	response := BaseResponse{
		Ok: false,
		Error: &BattleError{
			ErrorCode:    error_codes.ErrorCode(e.Code),
			ErrorMessage: fmt.Sprint(e.Message),
		},
	}
	if err2 := c.JSON(e.Code, response); err2 != nil {
		log.Error().Msg("error in converting HttpError to json: " + err2.Error())
	}
}
