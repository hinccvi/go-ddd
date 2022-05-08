package accesslog

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
)

func Handler(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Errors.Last() != nil {
			return
		}

		start := time.Now()

		ctx := c.Request.Context()
		ctx = log.WithRequest(ctx, c.Request)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		logger.With(ctx, "duration", time.Since(start).Milliseconds(), "status", http.StatusOK).
			Infof("%s %s %s %d %s", c.Request.Method, c.Request.URL.Path, c.Request.Proto, http.StatusOK, fmt.Sprintf("%dbytes", c.Request.ContentLength))
	}
}
