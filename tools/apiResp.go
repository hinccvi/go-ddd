package tools

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SuccessMsg = "success"
	CreatedMsg = "created"
)

func RespOkWithMsg[I any](c *gin.Context, msg string, i I) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  msg,
		"data": i,
	})
}
