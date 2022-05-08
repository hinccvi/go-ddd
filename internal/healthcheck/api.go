package healthcheck

import (
	"github.com/gin-gonic/gin"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/tools"
)

func RegisterHandlers(dg *gin.RouterGroup, version string) {
	dg.POST("/healthcheck", healthcheck(version))
}

func healthcheck(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tools.RespWithOK(c, "OK "+version)
	}
}
