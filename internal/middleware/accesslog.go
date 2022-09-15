package middleware

import (
	"fmt"
	"time"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func AccessLogHandler(logger log.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			fields := []zapcore.Field{
				zap.String("remote_ip", c.RealIP()),
				zap.String("latency", time.Since(start).String()),
				zap.String("host", req.Host),
				zap.String("request", fmt.Sprintf("%s %s", req.Method, req.RequestURI)),
				zap.Int("status", res.Status),
				zap.Int64("size", res.Size),
				zap.String("user_agent", req.UserAgent()),
				zap.String("request_id", id),
				zap.Error(err),
			}

			n := res.Status
			switch {
			case n >= 500:
				logger.Error("Server error", fields)
			case n >= 400:
				logger.Warn("Client error", fields)
			case n >= 300:
				logger.Info("Redirection", fields)
			default:
				logger.Info("Success", fields)
			}

			return nil
		}
	}
}
