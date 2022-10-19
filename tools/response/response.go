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

func Resp(c echo.Context, code int, msg message) error {
	statusCode := generateStatusCode(code)

	return c.JSON(statusCode, response{
		Code:    code,
		Message: string(msg),
	})
}

func RespWithData[I any](c echo.Context, code int, msg message, i I) error {
	statusCode := generateStatusCode(code)

	return c.JSON(statusCode, response{
		Code:    code,
		Message: string(msg),
		Data:    i,
	})
}
