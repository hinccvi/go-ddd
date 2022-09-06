package tools

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type message string

const (
	Success message = "success"
	Created message = "created"
	Updated message = "updated"
	Deleted message = "deleted"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RespOkWithData[I any](c *gin.Context, msg message, i I) {
	c.JSON(http.StatusOK, response{
		Code:    http.StatusOK,
		Message: string(msg),
		Data:    i,
	})
}

func RespOk(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, response{
		Code:    http.StatusOK,
		Message: msg,
	})
}
