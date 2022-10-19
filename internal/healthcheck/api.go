package healthcheck

import (
	"net/http"

	tools "github.com/hinccvi/Golang-Project-Structure-Conventional/tools/response"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(g *echo.Group, version string) {
	g.GET("/healthcheck", healthcheck(version))
}

func healthcheck(version string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return tools.RespWithData(c, http.StatusOK, tools.Success, "OK "+version)
	}
}
