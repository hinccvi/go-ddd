package tools

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	message string

	response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

const (
	Success message = "success"
	Created message = "created"
	Updated message = "updated"
	Deleted message = "deleted"
	Error   message = "error"
)

func generateStatusCode(code int) int {
	if code > http.StatusNetworkAuthenticationRequired {
		code = http.StatusBadRequest
	}

	return code
}

func BindValidate[T any](c echo.Context, t T) error {
	if err := c.Bind(t); err != nil {
		return err
	}

	return c.Validate(t)
}

func JSONRespOk(c echo.Context, i interface{}) error {
	return c.JSON(http.StatusOK, response{
		Code:    http.StatusOK,
		Message: string(Success),
		Data:    i,
	})
}

func JSONRespErr(c echo.Context, code int, i interface{}) error {
	statusCode := generateStatusCode(code)

	return c.JSON(statusCode, response{
		Code:    code,
		Message: string(Error),
		Data:    i,
	})
}
