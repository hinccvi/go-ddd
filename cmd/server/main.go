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
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	v1AuthController "github.com/hinccvi/go-ddd/internal/auth/controller/http/v1"
	authRepo "github.com/hinccvi/go-ddd/internal/auth/repository"
	authService "github.com/hinccvi/go-ddd/internal/auth/service"
	"github.com/hinccvi/go-ddd/internal/config"
	errs "github.com/hinccvi/go-ddd/internal/errors"
	hcController "github.com/hinccvi/go-ddd/internal/healthcheck/controller/http"
	m "github.com/hinccvi/go-ddd/internal/middleware"
	v1UserController "github.com/hinccvi/go-ddd/internal/user/controller/http/v1"
	userRepo "github.com/hinccvi/go-ddd/internal/user/repository"
	userService "github.com/hinccvi/go-ddd/internal/user/service"
	"github.com/hinccvi/go-ddd/pkg/db"
	"github.com/hinccvi/go-ddd/pkg/log"
	rds "github.com/hinccvi/go-ddd/pkg/redis"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	//nolint:gochecknoglobals // value of ldflags must be a package level variable
	Version = "1.0.0"

	//nolint:gochecknoglobals // environment flag that only used in main
	flagEnv = flag.String("env", "local", "environment")
)

const (
	gracefulTimeout   = 10 * time.Second
	readHeaderTimeout = 2 * time.Second
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
	db, err := db.Connect(ctx, &cfg)
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
		Handler:           buildHandler(logger, rds, db, &cfg),
		ReadHeaderTimeout: readHeaderTimeout,
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

	ctx, cancel := context.WithTimeout(ctx, gracefulTimeout)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		logger.Info(err)
	}

	logger.Info("Server exiting")
}

// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(logger log.Logger, rds redis.Client, db *sqlx.DB, cfg *config.Config) *echo.Echo {
	t := time.Duration(cfg.Context.Timeout) * time.Second

	e := echo.New()
	e.HTTPErrorHandler = m.NewHTTPErrorHandler(errs.GetStatusCodeMap()).Handler(logger)
	e.Validator = &m.CustomValidator{Validator: validator.New()}
	e.Use(buildMiddleware()...)

	authHandler := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &authService.JWTCustomClaims{},
		SigningKey: []byte(cfg.Jwt.AccessSigningKey),
	})

	dg := e.Group("")

	hcController.RegisterHandlers(
		dg,
		Version,
	)

	v1AuthController.RegisterHandlers(
		dg.Group("/v1"),
		authService.New(cfg, rds, authRepo.New(db, logger), logger, t),
		logger,
	)

	v1UserController.RegisterHandlers(
		dg.Group("/v1"),
		userService.New(rds, userRepo.New(db, logger), logger, t),
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

		middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Generator: func() string {
				return uuid.New().String()
			},
		}),

		// Api access log
		m.AccessLogHandler(logger),
	)

	return middlewares
}
