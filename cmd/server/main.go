package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1AuthPB "github.com/hinccvi/go-ddd/internal/auth/controller/grpc/v1"
	authRepo "github.com/hinccvi/go-ddd/internal/auth/repository"
	authService "github.com/hinccvi/go-ddd/internal/auth/service"
	"github.com/hinccvi/go-ddd/internal/config"
	v1UserPB "github.com/hinccvi/go-ddd/internal/user/controller/grpc/v1"
	userRepo "github.com/hinccvi/go-ddd/internal/user/repository"
	userService "github.com/hinccvi/go-ddd/internal/user/service"
	"github.com/hinccvi/go-ddd/pkg/db"
	"github.com/hinccvi/go-ddd/pkg/log"
	rds "github.com/hinccvi/go-ddd/pkg/redis"
	"github.com/hinccvi/go-ddd/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
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
		logger.Fatalf("fail to load app config: %v", err)
	}

	// connect to database
	db, err := db.Connect(ctx, &cfg)
	if err != nil {
		logger.Fatalf("fail to connect to db: %v", err)
	}

	// connect to redis
	rds, err := rds.Connect(ctx, &cfg)
	if err != nil {
		logger.Fatalf("fail to connect to redis: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.App.Port))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	if cfg.App.Cert == "" {
		cfg.App.Cert = data.Path("x509/server_cert.pem")
	}
	if cfg.App.Key == "" {
		cfg.App.Key = data.Path("x509/server_key.pem")
	}
	creds, err := credentials.NewServerTLSFromFile(cfg.App.Cert, cfg.App.Key)
	if err != nil {
		logger.Fatalf("failed to generate credentials: %v", err)
	}
	opts = []grpc.ServerOption{grpc.Creds(creds)}

	// timeout duration for each request
	t := time.Duration(cfg.Context.Timeout) * time.Second

	// register grpc server
	grpcServer := grpc.NewServer(opts...)

	proto.RegisterAuthServiceServer(
		grpcServer,
		v1AuthPB.RegisterHandlers(authService.New(&cfg, rds, authRepo.New(db, logger), logger, t), logger),
	)

	proto.RegisterUserServiceServer(
		grpcServer,
		v1UserPB.RegisterHandlers(userService.New(rds, userRepo.New(db, logger), logger, t), logger),
	)

	// seperate goroutine to listen on kill signal
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logger.Info("Server shutting down")

		grpcServer.GracefulStop()

		logger.Info("Server exiting")
	}()

	logger.Infof("grpc server listening on %v", lis.Addr())

	if err = grpcServer.Serve(lis); err != nil {
		logger.Fatalf("failed to serve grpc server: %v", err)
	}
}
