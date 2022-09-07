package tools

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type message string

const (
	MsgSuccess message = "success"
	MsgCreated message = "created"
	MsgUpdated message = "updated"
	MsgDeleted message = "deleted"
	MsgError   message = "error"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func generateStatusCode(code int) int {
	if code > 999 {
		code = http.StatusBadRequest
	}

	return code
}

func RespOkWithData[I any](c echo.Context, code int, msg message, i I) error {
	statusCode := generateStatusCode(code)

	return c.JSON(statusCode, response{
		Code:    code,
		Message: string(msg),
		Data:    i,
	})
}

func RespOk(c echo.Context, code int, msg message) error {
	statusCode := generateStatusCode(code)

	return c.JSON(statusCode, response{
		Code:    code,
		Message: string(msg),
	})
}
