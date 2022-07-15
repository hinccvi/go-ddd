package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/tools"
)

func RegisterHandlers(dg *gin.RouterGroup, service Service, logger log.Logger) {
	dg.POST("/login", login(service, logger))
}

func login(service Service, logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(err)
			return
		}

		token, err := service.Login(c.Request.Context(), req.Name, req.Password)
		if err != nil {
			c.Error(err)
			return
		}

		tools.RespOkWithData(c, tools.SuccessMsg, struct {
			Token string `json:"token"`
		}{token})
	}
}
