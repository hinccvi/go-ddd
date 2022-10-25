package v1

import (
	"net/http"
	"strings"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/auth/service"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/tools"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(g *echo.Group, service service.Service, logger log.Logger) {
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
	service service.Service
}

func (r resource) login(c echo.Context) error {
	var req service.LoginRequest
	if err := tools.BindValidate(c, &req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	res, err := r.service.Login(ctx, req)
	if err != nil {
		return err
	}

	return tools.JSONRespWithData(c, http.StatusOK, tools.Success, res)
}

func (r resource) refresh(c echo.Context) error {
	var req service.RefreshTokenRequest
	if err := tools.BindValidate(c, &req); err != nil {
		return err
	}

	tokenString := c.Request().Header.Get("Authorization")
	accessTokenArr := strings.Split(tokenString, " ")
	req.AccessToken = accessTokenArr[1]

	if len(accessTokenArr) != constants.JWTpart {
		return constants.ErrInvalidJwt
	}

	ctx := c.Request().Context()
	res, err := r.service.Refresh(ctx, req)
	if err != nil {
		return err
	}

	return tools.JSONRespWithData(c, http.StatusOK, tools.Success, res)
}
