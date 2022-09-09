package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/auth"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/config"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/healthcheck"
	m "github.com/hinccvi/Golang-Project-Structure-Conventional/internal/middleware"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/user"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/accesslog"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/db"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	rds "github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/redis"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/tools"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var Version = "1.0.0"

var flagMode = flag.String("mode", "local", "environment")

func main() {
	flag.Parse()

	// create root logger tagged with server version
	logger := log.New(*flagMode, zap.ErrorLevel).With(context.TODO(), "version", Version)

	// load application configurations
	cfg, err := config.Load(*flagMode)
	if err != nil {
		logger.Fatal(err)
	}

	// connect to database
	dbx, err := db.Connect(*flagMode, &cfg)
	if err != nil {
		logger.Fatal(err)
	}

	// connect to redis
	rds, err := rds.Connect(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: buildHandler(*flagMode, &logger, rds, &dbx, &cfg),
	}

	logger.Infof("Server listening on %s", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Server exiting")
}

// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(mode string, logger *log.Logger, rds *redis.Client, dbx *models.DBTX, cfg *config.Config) *echo.Echo {
	e := echo.New()

	e.HTTPErrorHandler = m.NewHttpErrorHandler(constants.ErrorStatusCodeMaps).Handler(*logger)

	e.Use(buildMiddleware()...)

	e.Validator = &m.CustomValidator{Validator: validator.New()}

	authHandler := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &auth.JwtCustomClaims{},
		SigningKey: []byte(cfg.Jwt.AccessSigningKey),
	})

	defaultGroup := e.Group("")

	healthcheck.RegisterHandlers(
		defaultGroup,
		Version,
	)

	auth.RegisterHandlers(
		defaultGroup,
		auth.NewService(cfg, auth.NewRepository(dbx, *logger), *logger),
		*logger,
	)

	user.RegisterHandlers(
		defaultGroup,
		user.NewService(rds, user.NewRepository(dbx, *logger), *logger),
		*logger,
		authHandler,
	)

	return e
}

// buildMiddleware sets up the middlewre logic and builds a handler.
func buildMiddleware() []echo.MiddlewareFunc {
	var middlewares []echo.MiddlewareFunc
	logger := log.New(*flagMode, zap.InfoLevel).With(context.TODO(), "version", Version)

	middlewares = append(middlewares,

		// Recover
		middleware.Recover(),

		// Api access logs
		accesslog.Handler(logger),

		// X-Request-ID
		middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Generator: func() string {
				return tools.GenerateUUIDv4().String()
			},
		}),
	)

	return middlewares
}
