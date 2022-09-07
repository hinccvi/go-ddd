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

	"github.com/go-redis/redis/v8"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/auth"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/config"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/errors"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/healthcheck"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/user"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/accesslog"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/db"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	rds "github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/swaggo/echo-swagger/example/docs"
	"gorm.io/gorm"
)

var Version = "1.0.0"

var flagMode = flag.String("mode", "local", "environment")

// @title Swagger Example API
// @version 1.1
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8022
// @BasePath /v1
func main() {
	flag.Parse()

	// create root logger tagged with server version
	logger := log.New(*flagMode).With(context.TODO(), "version", Version)

	// load application configurations
	cfg, err := config.Load(*flagMode)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(1)
	}

	// connect to database
	dbx, err := db.Connect(*flagMode, cfg)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	// connect to redis
	rds, err := rds.Connect(cfg)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: buildHandler(*flagMode, &logger, rds, dbx, &cfg),
	}

	logger.Infof("Server listening on %s", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	logger.Info("Server exiting")
}

// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(mode string, logger *log.Logger, rds *redis.Client, dbx *gorm.DB, cfg *config.Config) *echo.Echo {
	e := echo.New()

	e.HTTPErrorHandler = errors.NewHttpErrorHandler(constants.ErrorStatusCodeMaps).Handler(*logger)

	e.Use(
		accesslog.Handler(*logger),
	)

	authHandler := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &auth.JwtCustomClaims{},
		SigningKey: []byte(cfg.Jwt.AccessSigningKey),
	})

	defaultGroup := e.Group("")

	defaultGroup.GET("/swagger/*", echoSwagger.WrapHandler)

	healthcheck.RegisterHandlers(
		defaultGroup,
		Version,
		authHandler,
	)

	auth.RegisterHandlers(
		defaultGroup,
		auth.NewService(cfg.Jwt.AccessSigningKey, cfg.Jwt.AccessExpiration, auth.NewRepository(dbx, *logger), *logger),
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
