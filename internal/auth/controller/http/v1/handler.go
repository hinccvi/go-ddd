package v1

import (
	"strings"

	"github.com/hinccvi/go-ddd/internal/auth/service"
	errs "github.com/hinccvi/go-ddd/internal/errors"
	"github.com/hinccvi/go-ddd/pkg/log"
	"github.com/hinccvi/go-ddd/tools"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(g *echo.Group, service service.Service, logger log.Logger) {
	r := &resource{logger, service}

	auth := g.Group("/auth")
	{
		auth.POST("/login", r.login)
		auth.POST("/refresh", r.refresh)
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

	return tools.JSONRespOk(c, res)
}

func (r resource) refresh(c echo.Context) error {
	var req service.RefreshTokenRequest
	if err := tools.BindValidate(c, &req); err != nil {
		return err
	}

	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return errs.ErrInvalidJwt
	}

	accessTokenArr := strings.Split(tokenString, " ")
	req.AccessToken = accessTokenArr[1]

	if len(accessTokenArr) != service.JWTBearerFormat {
		return errs.ErrInvalidJwt
	}

	ctx := c.Request().Context()
	res, err := r.service.Refresh(ctx, req)
	if err != nil {
		return err
	}

	return tools.JSONRespOk(c, &res)
}
