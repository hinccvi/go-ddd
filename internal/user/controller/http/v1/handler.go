package v1

import (
	"context"
	"net/http"

	"github.com/hinccvi/go-ddd/internal/entity"
	"github.com/hinccvi/go-ddd/internal/user/service"
	"github.com/hinccvi/go-ddd/pkg/log"
	"github.com/hinccvi/go-ddd/tools"
	"github.com/labstack/echo/v4"
)

type resource struct {
	logger  log.Logger
	service service.Service
}

func RegisterHandlers(g *echo.Group, service service.Service, logger log.Logger, authHandler echo.MiddlewareFunc) {
	r := &resource{logger, service}

	user := g.Group("/user")
	{
		user.GET("/:id", r.Get)
		user.GET("/list", r.Query)
		user.POST("", r.Create)

		user.PATCH("", r.Update, authHandler)
		user.DELETE("/:id", r.Delete, authHandler)
	}
}

func (r resource) Get(c echo.Context) error {
	var req service.GetUserRequest
	if err := tools.BindValidate(c, &req); err != nil {
		return err
	}

	user, err := r.service.Get(c.Request().Context(), *req.ID)
	if err != nil {
		return err
	}

	return tools.JSON(c, http.StatusOK, tools.Success, user)
}

func (r resource) Query(c echo.Context) error {
	var req service.QueryUserRequest
	if err := tools.BindValidate(c, &req); err != nil {
		return err
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 10
	}

	list, total, err := r.service.Query(c.Request().Context(), req.Page, req.Size)
	if err != nil {
		return err
	}

	return tools.JSON(
		c,
		http.StatusOK,
		tools.Success,
		struct {
			List  []entity.User `json:"list"`
			Total int64         `json:"total"`
		}{
			List:  list,
			Total: total,
		})
}

func (r resource) Create(c echo.Context) error {
	var req service.CreateUserRequest
	if err := tools.BindValidate(c, &req); err != nil {
		return err
	}

	u := entity.User{
		Username: req.Username,
		Password: req.Password,
	}
	if err := r.service.Create(context.TODO(), u); err != nil {
		return err
	}

	return tools.JSON(c, http.StatusOK, tools.Created, nil)
}

func (r resource) Update(c echo.Context) error {
	var req service.UpdateUserRequest
	if err := tools.BindValidate(c, &req); err != nil {
		return err
	}

	u := entity.User{
		ID:       req.ID,
		Username: req.Username,
		Password: req.Password,
	}
	if err := r.service.Update(context.TODO(), u); err != nil {
		return err
	}

	return tools.JSON(c, http.StatusOK, tools.Updated, nil)
}

func (r resource) Delete(c echo.Context) error {
	var req service.DeleteUserRequest
	if err := tools.BindValidate(c, &req); err != nil {
		return err
	}

	if err := r.service.Delete(context.TODO(), *req.ID); err != nil {
		return err
	}

	return tools.JSON(c, http.StatusOK, tools.Deleted, nil)
}
