package user

import (
	"github.com/gin-gonic/gin"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/tools"
)

func RegisterHandlers(dg *gin.RouterGroup, service Service, authHandler gin.HandlerFunc, logger log.Logger) {
	r := &resource{logger, service}

	v1 := dg.Group("v1")
	{
		v1.GET("/get", r.Get)
		v1.GET("/query", r.Query)
		v1.GET("/count", r.Count)
		v1.POST("/create", r.Create)

		v1.Use(authHandler)

		v1.PATCH("/update", r.Update)
		v1.DELETE("/delete", r.Delete)
	}
}

type resource struct {
	logger  log.Logger
	service Service
}

func (r resource) Get(c *gin.Context) {
	var req getOrDeleteUserRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	user, err := r.service.Get(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	tools.RespOkWithData(c, tools.SuccessMsg, user)
}

func (r resource) Query(c *gin.Context) {
	var req queryUserRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	users, err := r.service.Query(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	tools.RespOkWithData(c, tools.SuccessMsg, users)
}

func (r resource) Count(c *gin.Context) {
	total, err := r.service.Count(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	tools.RespOkWithData(c, tools.SuccessMsg, struct {
		Total int64 `json:"total"`
	}{
		Total: total,
	})
}

func (r resource) Create(c *gin.Context) {
	var req createUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	user, err := r.service.Create(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	tools.RespOkWithData(c, tools.CreatedMsg, user)
}

func (r resource) Update(c *gin.Context) {
	var req updateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	user, err := r.service.Update(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	tools.RespOkWithData(c, tools.SuccessMsg, user)
}

func (r resource) Delete(c *gin.Context) {
	var req getOrDeleteUserRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}

	user, err := r.service.Delete(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	tools.RespOkWithData(c, tools.SuccessMsg, user)
}
