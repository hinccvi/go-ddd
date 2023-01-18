package v1

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/hinccvi/go-ddd/pkg/log"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	authService "github.com/hinccvi/go-ddd/internal/auth/service"
	"github.com/hinccvi/go-ddd/internal/config"
	"github.com/hinccvi/go-ddd/internal/mocks"
	"github.com/hinccvi/go-ddd/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	authRepo "github.com/hinccvi/go-ddd/internal/auth/repository"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var lis *bufconn.Listener

const bufSize = 1024 * 1024

func TestMain(m *testing.M) {
	l, _ := log.NewForTest()
	logger := log.NewWithZap(l)

	db, _, err := sqlmock.New()
	if err != nil {
		logger.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbx := sqlx.NewDb(db, "pgx")
	defer db.Close()

	var cfg config.Config
	cfg.Jwt.AccessExpiration = int(5 * time.Minute)
	cfg.Jwt.AccessSigningKey = "secret1"
	cfg.Jwt.RefreshExpiration = int(7 * time.Hour * 24)
	cfg.Jwt.RefreshSigningKey = "secret2"

	mnr, _ := miniredis.Run()
	defer mnr.Close()

	rds, err := mocks.Redis(mnr.Addr())

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	proto.RegisterAuthServiceServer(
		s,
		RegisterHandlers(authService.New(&cfg, rds, authRepo.New(dbx, logger), logger, 5*time.Second), logger),
	)

	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Fatalf("server exited with error: %v", err)
		}
	}()

	defer s.GracefulStop()

	m.Run()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != err {
		t.Fatalf("failed to dial grpc server: %v", err)
	}
	defer conn.Close()
	c := proto.NewAuthServiceClient(conn)

	t.Run("fail: empty param", func(t *testing.T) {
		_, err := c.Login(ctx, &proto.LoginRequest{})
		assert.NoError(t, err)
	})
}
