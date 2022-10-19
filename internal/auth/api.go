package auth

import (
	"net/http"
	"strings"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	tools "github.com/hinccvi/Golang-Project-Structure-Conventional/tools/response"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(g *echo.Group, service Service, logger log.Logger) {
	r := &resource{logger, service}

	v1 := g.Group("v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", r.login)
			auth.POST("/refresh", r.refresh)
		}
	}
}

type resource struct {
	logger  log.Logger
	service Service
}

func (r resource) login(c echo.Context) error {
	var req loginRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	res, err := r.service.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return err
	}

	return tools.RespWithData(c, http.StatusOK, tools.Success, res)
}

func (r resource) refresh(c echo.Context) error {
	var req refreshRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	tokenString := c.Request().Header.Get("Authorization")
	accessToken := strings.Split(tokenString, " ")

	res, err := r.service.Refresh(c.Request().Context(), accessToken[1], req.RefreshToken)
	if err != nil {
		return err
	}

	return tools.RespWithData(c, http.StatusOK, tools.Success, res)
}
