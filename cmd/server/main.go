package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v9"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/auth"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/config"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/healthcheck"
	m "github.com/hinccvi/Golang-Project-Structure-Conventional/internal/middleware"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/models"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/user"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/db"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	rds "github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/redis"
	uTools "github.com/hinccvi/Golang-Project-Structure-Conventional/tools/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	//nolint:gochecknoglobals // value of ldflags must be a package level variable
	Version = "1.0.0"

	//nolint:gochecknoglobals // environment flag that only used in main
	flagEnv = flag.String("env", "local", "environment")
)

func main() {
	flag.Parse()

	// create root context
	ctx := context.Background()

	// create root logger tagged with server version
	logger := log.NewWithZap(log.New(*flagEnv, log.ErrorLog)).With(ctx, "version", Version)

	// load application configurations
	cfg, err := config.Load(*flagEnv)
	if err != nil {
		logger.Fatal(err)
	}

	// connect to database
	dbx, err := db.Connect(&cfg, log.New(*flagEnv, log.SQLLog))
	if err != nil {
		logger.Fatal(err)
	}

	// connect to redis
	rds, err := rds.Connect(ctx, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.App.Port),
		Handler:           buildHandler(logger, rds, dbx, &cfg),
		ReadHeaderTimeout: constants.RequestTimeout,
	}

	logger.Infof("Server listening on %s", server.Addr)

	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down")

	ctx, cancel := context.WithTimeout(ctx, constants.ContextTimeoutDuration)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		logger.Info(err)
	}

	logger.Info("Server exiting")
}

// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(logger log.Logger, rds redis.Client, dbx models.DBTX, cfg *config.Config) *echo.Echo {
	e := echo.New()

	e.HTTPErrorHandler = m.NewHTTPErrorHandler(constants.ErrorStatusCodeMaps).Handler(logger)

	e.Use(buildMiddleware()...)

	e.Validator = &m.CustomValidator{Validator: validator.New()}

	authHandler := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &constants.JWTCustomClaims{},
		SigningKey: []byte(cfg.Jwt.AccessSigningKey),
	})

	defaultGroup := e.Group("")

	healthcheck.RegisterHandlers(
		defaultGroup,
		Version,
	)

	auth.RegisterHandlers(
		defaultGroup,
		auth.NewService(cfg, rds, auth.NewRepository(dbx, logger), logger),
		logger,
	)

	user.RegisterHandlers(
		defaultGroup,
		user.NewService(rds, user.NewRepository(dbx, logger), logger),
		logger,
		authHandler,
	)

	return e
}

// buildMiddleware sets up the middlewre logic and builds a handler.
func buildMiddleware() []echo.MiddlewareFunc {
	var middlewares []echo.MiddlewareFunc
	logger := log.NewWithZap(log.New(*flagEnv, log.AccessLog)).With(context.TODO(), "version", Version)

	middlewares = append(middlewares,

		// Echo built-in middleware
		middleware.Recover(),

		middleware.Secure(),

		middleware.TimeoutWithConfig(middleware.TimeoutConfig{
			Timeout:      constants.RequestTimeout,
			ErrorMessage: constants.MsgRequestTimeout,
		}),

		middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Generator: func() string {
				u, err := uTools.GenerateUUIDv4()
				for err != nil {
					u, err = uTools.GenerateUUIDv4()
				}

				return u.String()
			},
		}),

		// Api access log
		m.AccessLogHandler(logger),
	)

	return middlewares
}
