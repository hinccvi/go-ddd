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

func RespOkWithData[I any](c echo.Context, code int, msg message, i I) error {
	status := http.StatusBadRequest
	if code < 999 {
		status = code
	}

	return c.JSON(status, response{
		Code:    code,
		Message: string(msg),
		Data:    i,
	})
}

func RespOk(c echo.Context, code int, msg message) error {
	status := http.StatusBadRequest
	if code < 999 {
		status = code
	}

	return c.JSON(status, response{
		Code:    code,
		Message: string(msg),
	})
}
