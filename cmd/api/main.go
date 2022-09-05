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

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/auth"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/config"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/errors"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/healthcheck"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/user"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/accesslog"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/db"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	rds "github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/redis"
	"github.com/mattn/go-colorable"
	"gorm.io/gorm"
)

var Version = "1.0.0"

var flagMode = flag.String("mode", "local", "environment")

func main() {
	flag.Parse()

	// create root logger tagged with server version
	logger := log.New(*flagMode).With(context.TODO(), "version", Version)

	// load application configurations
	cfg, err := config.Load(*flagMode)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	// connect to database
	dbx, err := db.Connect(*flagMode, cfg)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	// connect to redis
	rds, err := rds.Connect(cfg)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.AppConfig.Port),
		Handler: buildHandler(*flagMode, logger, rds, dbx, cfg),
	}

	logger.Infof("Server listening on %s", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(err)
			os.Exit(-1)
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
		os.Exit(-1)
	}

	logger.Info("Server exiting")
}

// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(mode string, logger log.Logger, rds *redis.Client, dbx *gorm.DB, cfg config.Config) *gin.Engine {
	if mode == "local" {
		gin.ForceConsoleColor()
		gin.DefaultWriter = colorable.NewColorableStdout()
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.Default()
	e.Use(
		accesslog.Handler(logger),
		errors.Handler(logger),
	)
	e.NoRoute(func(c *gin.Context) {
		c.Error(errors.NotFound("resource not found"))
	})

	authHandler := auth.Handler(cfg.JwtConfig.AccessSigningKey)

	defaultRouterGroup := e.Group("")

	healthcheck.RegisterHandlers(
		defaultRouterGroup,
		Version,
	)

	auth.RegisterHandlers(
		defaultRouterGroup,
		auth.NewService(cfg.JwtConfig.AccessSigningKey, cfg.JwtConfig.AccessExpiration, auth.NewRepository(dbx, logger), logger),
		logger,
	)

	user.RegisterHandlers(
		defaultRouterGroup,
		user.NewService(rds, user.NewRepository(dbx, logger), logger),
		authHandler,
		logger,
	)

	return e
}
