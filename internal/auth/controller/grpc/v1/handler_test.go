package v1

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	logger "github.com/hinccvi/go-ddd/pkg/log"
	"github.com/hinccvi/go-ddd/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	authService "github.com/hinccvi/go-ddd/internal/auth/service"
	"github.com/hinccvi/go-ddd/internal/config"
	"github.com/hinccvi/go-ddd/internal/entity"
	"github.com/hinccvi/go-ddd/internal/mocks"
	"github.com/hinccvi/go-ddd/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	lis      *bufconn.Listener
	client   pb.AuthServiceClient
	password string
)

const bufSize = 1024 * 1024

func TestMain(m *testing.M) {
	s := serverSetup()
	defer s.GracefulStop()
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("server exited with error: %v", err)
		}
	}()

	conn := clientSetup()
	defer conn.Close()
	client = pb.NewAuthServiceClient(conn)

	m.Run()
}

func clientSetup() *grpc.ClientConn {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != err {
		log.Fatalf("failed to dial grpc server: %v", err)
	}
	return conn
}

func serverSetup() *grpc.Server {
	l, _ := logger.NewForTest()
	logger := logger.NewWithZap(l)

	var cfg config.Config
	cfg.Jwt.AccessExpiration = int(10 * time.Second)
	cfg.Jwt.AccessSigningKey = "secret1"
	cfg.Jwt.RefreshExpiration = int(7 * time.Hour * 24)
	cfg.Jwt.RefreshSigningKey = "secret2"

	mnr, _ := miniredis.Run()

	rds, _ := mocks.Redis(mnr.Addr())

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	password, _ = tools.Bcrypt("secret")
	mockGetUserByUsername := entity.User{
		ID:       uuid.New(),
		Username: "user",
		Password: password,
	}

	var repo mocks.AuthRepository
	repo.On("GetUserByUsername", mock.Anything, "user").Return(mockGetUserByUsername, nil)

	pb.RegisterAuthServiceServer(
		s,
		RegisterHandlers(authService.New(&cfg, rds, &repo, logger, 5*time.Second), logger),
	)

	return s
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestLogin(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		reply, err := client.Login(ctx, &pb.LoginRequest{Username: "user", Password: "secret"})
		assert.NoError(t, err)
		assert.NotNil(t, reply)
	})

	t.Run("fail: incorrect credentials", func(t *testing.T) {
		_, err := client.Login(ctx, &pb.LoginRequest{Username: "user", Password: "password"})
		assert.Error(t, err)
	})

	t.Run("fail: empty param", func(t *testing.T) {
		_, err := client.Login(ctx, &pb.LoginRequest{})
		assert.Error(t, err)
	})
}

func TestRefresh(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		loginReply, err := client.Login(ctx, &pb.LoginRequest{Username: "user", Password: "secret"})
		assert.NoError(t, err)
		assert.NotNil(t, loginReply)

		header := metadata.New(map[string]string{"Authorization": loginReply.AccessToken})
		ctx = metadata.NewOutgoingContext(ctx, header)

		reply, err := client.Refresh(ctx, &pb.RefreshRequest{RefreshToken: loginReply.RefreshToken})
		assert.NoError(t, err)
		assert.NotNil(t, reply)
	})

	t.Run("fail: empty param", func(t *testing.T) {
		_, err := client.Refresh(ctx, &pb.RefreshRequest{})
		assert.Error(t, err)
	})
}
