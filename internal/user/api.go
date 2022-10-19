package user

import (
	"context"
	"net/http"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
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

	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, ctxID, req.ID)

	user, err := r.service.Get(ctx)
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

	args := models.ListUserParams{
		Limit:  limit,
		Offset: offset,
	}

	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, ctxListUser, args)

	list, err := r.service.Query(ctx)
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

	args := models.CreateUserParams{
		Username: req.Username,
		Password: req.Password,
	}

	ctx := context.TODO()
	ctx = context.WithValue(ctx, ctxCreateUser, args)

	user, err := r.service.Create(ctx)
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

	args := models.UpdateUserParams{
		ID:       req.ID,
		Username: req.Username,
		Password: req.Password,
	}

	ctx := context.TODO()
	ctx = context.WithValue(ctx, ctxUpdateUser, args)

	user, err := r.service.Update(ctx)
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

	ctx := context.TODO()
	ctx = context.WithValue(ctx, ctxID, req.ID)

	user, err := r.service.Delete(ctx)
	if err != nil {
		return err
	}

	return tools.RespWithData(c, http.StatusOK, tools.Deleted, user)
}
