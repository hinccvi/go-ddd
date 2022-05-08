package test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/errors"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/accesslog"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
)

// MockRoutingContext creates a routing.Conext for testing handlers.
func MockRoutingContext(req *http.Request) (*gin.Context, *httptest.ResponseRecorder) {
	res := httptest.NewRecorder()
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	ctx, _ := gin.CreateTestContext(res)
	return ctx, res
}

// MockRouter creates a routing.Router for testing APIs.
func MockRouter(logger log.Logger) *gin.Engine {
	e := gin.Default()
	e.Use(
		accesslog.Handler(logger),
		errors.Handler(logger),
	)
	return e
}
