package user

import (
	"context"
	"net/http"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/tools"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(g *echo.Group, service Service, logger log.Logger, authHandler echo.MiddlewareFunc) {
	r := &resource{logger, service}

	v1 := g.Group("v1")
	{
		user := v1.Group("/user")
		{
			user.GET("", r.Get)
			user.GET("/list", r.Query)
			user.GET("/count", r.Count)
			user.POST("", r.Create)

			// user.Use(authHandler)

			user.PATCH("", r.Update)
			user.DELETE("", r.Delete)
		}
	}
}

type resource struct {
	logger  log.Logger
	service Service
}

func (r resource) Get(c echo.Context) error {
	var req getOrDeleteUserRequest

	if err := tools.Validator(c, &req); err != nil {
		return err
	}

	user, err := r.service.Get(c.Request().Context(), req.Id)
	if err != nil {
		return err
	}

	return tools.RespOkWithData(c, http.StatusOK, tools.MsgSuccess, user)
}

func (r resource) Query(c echo.Context) error {
	var req queryUserRequest

	if err := tools.Validator(c, &req); err != nil {
		return err
	}

	var (
		limit  int32
		offset int32
	)
	if req.Limit > 0 {
		limit = int32(req.Limit)
	}

	if req.Offset != nil && *req.Offset > 0 {
		offset = int32(*req.Offset)
	}

	arg := &models.ListUserParams{
		Limit:  limit,
		Offset: offset,
	}

	users, err := r.service.Query(c.Request().Context(), arg)
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

	if err := tools.Validator(c, &req); err != nil {
		return err
	}

	arg := &models.CreateUserParams{
		Username: req.Username,
		Password: req.Password,
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), constants.RequestTimeoutDuration)
	defer cancel()

	user, err := r.service.Create(ctx, arg)
	if err != nil {
		return err
	}

	return tools.RespOkWithData(c, http.StatusOK, tools.MsgCreated, user)
}

func (r resource) Update(c echo.Context) error {
	var req updateUserRequest

	if err := tools.Validator(c, &req); err != nil {
		return err
	}

	arg := &models.UpdateUserParams{
		ID:       req.Id,
		Username: req.Username,
		Password: req.Password,
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), constants.RequestTimeoutDuration)
	defer cancel()

	user, err := r.service.Update(ctx, arg)
	if err != nil {
		return err
	}

	return tools.RespOkWithData(c, http.StatusOK, tools.MsgSuccess, user)
}

func (r resource) Delete(c echo.Context) error {
	var req getOrDeleteUserRequest

	if err := tools.Validator(c, &req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), constants.RequestTimeoutDuration)
	defer cancel()

	user, err := r.service.Delete(ctx, req.Id)
	if err != nil {
		return err
	}

	return tools.RespOkWithData(c, http.StatusOK, tools.MsgSuccess, user)
}
