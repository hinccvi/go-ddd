package user

import (
	"context"
	"net/http"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/model"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	tools "github.com/hinccvi/Golang-Project-Structure-Conventional/tools/response"
	"github.com/labstack/echo/v4"
)

type resource struct {
	logger  log.Logger
	service Service
}

func RegisterHandlers(g *echo.Group, service Service, logger log.Logger, authHandler echo.MiddlewareFunc) {
	r := &resource{logger, service}

	v1 := g.Group("v1")
	{
		user := v1.Group("/user")
		{
			user.GET("/:id", r.Get)
			user.GET("/list", r.Query)
			user.POST("", r.Create)

			user.PATCH("", r.Update, authHandler)
			user.DELETE("/:id", r.Delete, authHandler)
		}
	}
}

func (r resource) Get(c echo.Context) error {
	var req getUserRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	user, err := r.service.Get(c.Request().Context(), *req.ID)
	if err != nil {
		return err
	}

	return tools.RespWithData(c, http.StatusOK, tools.Success, user)
}

func (r resource) Query(c echo.Context) error {
	var req queryUserRequest

	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	var (
		limit  int32 = 10
		offset int32
	)

	if req.Limit > 0 {
		limit = req.Limit
	}

	if req.Offset != nil && *req.Offset > 0 {
		offset = *req.Offset
	}

	args := model.ListUserParams{
		Limit:  limit,
		Offset: offset,
	}

	list, err := r.service.Query(c.Request().Context(), args)
	if err != nil {
		return err
	}

	return tools.RespWithData(c, http.StatusOK, tools.Success, list)
}

func (r resource) Create(c echo.Context) error {
	var req createUserRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	args := model.CreateUserParams{
		Username: req.Username,
		Password: req.Password,
	}

	user, err := r.service.Create(context.TODO(), args)
	if err != nil {
		return err
	}

	return tools.RespWithData(c, http.StatusOK, tools.Created, user)
}

func (r resource) Update(c echo.Context) error {
	var req updateUserRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	args := model.UpdateUserParams{
		ID:       req.ID,
		Username: req.Username,
		Password: req.Password,
	}

	user, err := r.service.Update(context.TODO(), args)
	if err != nil {
		return err
	}

	return tools.RespWithData(c, http.StatusOK, tools.Updated, user)
}

func (r resource) Delete(c echo.Context) error {
	var req deleteUserRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	user, err := r.service.Delete(context.TODO(), *req.ID)
	if err != nil {
		return err
	}

	return tools.RespWithData(c, http.StatusOK, tools.Deleted, user)
}
