package tools

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespWithOK[I any](c *gin.Context, i I) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "Success",
		"data": i,
	})
}

func RespWithCreated[I any](c *gin.Context, i I) {
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusCreated,
		"msg":  "Created",
		"data": i,
	})
}
