package handler

import (
	"errors"

	pkgerr "github.com/ariefsibuea/flight-aggregator/internal/pkg/errors"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

type Error struct {
	Message string `json:"message"`
}

func Success(c echo.Context, statusCode int, data any) error {
	return c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

func Fail(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, Response{
		Success: false,
		Error:   &Error{Message: message},
	})
}

func ErrorHandler(err error, c echo.Context) {
	var (
		code    int
		message string
	)

	var apiErr *pkgerr.APIError
	switch {
	case errors.As(err, &apiErr):
		code = apiErr.Code()
		message = apiErr.Error()
	default:
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message.(string)
		} else {
			code = pkgerr.GetErrorCode(err)
			message = err.Error()
		}
	}

	c.JSON(code, Response{
		Success: false,
		Error:   &Error{Message: message},
	})
}
