package healthcheck

import (
	"net/http"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/tools"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(g *echo.Group, version string) {
	g.GET("/healthcheck", healthcheck(version))
}

func healthcheck(version string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return tools.RespOkWithData(c, http.StatusOK, tools.MsgSuccess, "OK "+version)
	}
}
