package user

import (
	"net/http"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/tools"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(g *echo.Group, service Service, logger log.Logger, authHandler echo.MiddlewareFunc) {
	r := &resource{logger, service}

	v1 := g.Group("v1")
	{
		v1.GET("/get", r.Get)
		v1.GET("/query", r.Query)
		v1.GET("/count", r.Count)
		v1.POST("/create", r.Create)

		// v1.Use(authHandler)

		v1.PATCH("/update", r.Update)
		v1.DELETE("/delete", r.Delete)
	}
}

type resource struct {
	logger  log.Logger
	service Service
}

func (r resource) Get(c echo.Context) error {
	var req getOrDeleteUserRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	user, err := r.service.Get(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return tools.RespOkWithData(c, http.StatusOK, tools.MsgSuccess, user)
}

func (r resource) Query(c echo.Context) error {
	var req queryUserRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	users, err := r.service.Query(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return tools.RespOkWithData(c, http.StatusOK, tools.MsgSuccess, users)
}

func (r resource) Count(c echo.Context) error {
	total, err := r.service.Count(c.Request().Context())
	if err != nil {
		return err
	}

	return tools.RespOkWithData(c, http.StatusOK, tools.MsgSuccess, struct {
		Total int64 `json:"total"`
	}{
		Total: total,
	})
}

func (r resource) Create(c echo.Context) error {
	var req createUserRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	user, err := r.service.Create(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return tools.RespOkWithData(c, http.StatusOK, tools.MsgCreated, user)
}

func (r resource) Update(c echo.Context) error {
	var req updateUserRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	user, err := r.service.Update(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return tools.RespOkWithData(c, http.StatusOK, tools.MsgSuccess, user)
}

func (r resource) Delete(c echo.Context) error {
	var req getOrDeleteUserRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	user, err := r.service.Delete(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return tools.RespOkWithData(c, http.StatusOK, tools.MsgSuccess, user)
}
