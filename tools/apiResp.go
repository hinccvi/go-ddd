package tools

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SuccessMsg = "success"
	CreatedMsg = "created"
)

type response struct {
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func RespOkWithData[I any](c *gin.Context, msg string, i I) {
	c.JSON(http.StatusOK, response{
		Message: msg,
		Data:    i,
	})
}

func RespOk(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, response{
		Message: msg,
	})
}
