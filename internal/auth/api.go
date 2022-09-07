package auth

import (
	"net/http"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/tools"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(g *echo.Group, service Service, logger log.Logger) {
	g.POST("/login", login(service, logger))
}

func login(service Service, logger log.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req loginRequest

		if err := c.Bind(&req); err != nil {
			return err
		}

		token, err := service.Login(c.Request().Context(), req.Name, req.Password)
		if err != nil {
			return err
		}

		return tools.RespOkWithData(c, http.StatusOK, tools.MsgSuccess, struct {
			Token string `json:"token"`
		}{token})
	}
}
