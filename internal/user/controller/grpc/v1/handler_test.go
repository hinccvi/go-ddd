package v1

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/google/uuid"
	"github.com/hinccvi/go-ddd/internal/entity"
	"github.com/hinccvi/go-ddd/internal/mocks"
	userService "github.com/hinccvi/go-ddd/internal/user/service"
	logger "github.com/hinccvi/go-ddd/pkg/log"
	"github.com/hinccvi/go-ddd/proto/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var (
	lis    *bufconn.Listener
	client pb.UserServiceClient
	id1    uuid.UUID = uuid.New()
	id2    uuid.UUID = uuid.New()
)

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
	client = pb.NewUserServiceClient(conn)

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

	mnr, _ := miniredis.Run()

	rds, _ := mocks.Redis(mnr.Addr())

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	repo := mocks.UserRepository{
		Items: []entity.User{
			{
				ID:        id1,
				Username:  "user1",
				CreatedAt: time.Now(),
			},
			{
				ID:        id2,
				Username:  "user2",
				CreatedAt: time.Now(),
			},
		},
	}

	pb.RegisterUserServiceServer(
		s,
		RegisterHandlers(userService.New(rds, &repo, logger, 5*time.Second), logger),
	)

	return s
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGetUser(t *testing.T) {
	ctx := context.TODO()

	t.Run("sucess", func(t *testing.T) {
		reply, err := client.GetUser(ctx, &pb.GetUserRequest{Id: id1.String()})
		assert.NoError(t, err)
		assert.NotNil(t, reply)
	})

	t.Run("fail: invalid uuid", func(t *testing.T) {
		_, err := client.GetUser(ctx, &pb.GetUserRequest{Id: "xxx"})
		assert.Error(t, err)
	})

	t.Run("fail: id not found", func(t *testing.T) {
		_, err := client.GetUser(ctx, &pb.GetUserRequest{Id: uuid.NewString()})
		assert.Error(t, err)
	})
}

func TestQueryUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		reply, err := client.QueryUser(ctx, &pb.QueryUserRequest{Page: 1, Size: 10})
		assert.NoError(t, err)
		if assert.NotNil(t, reply) {
			assert.Len(t, reply.Users, 2)
			assert.EqualValues(t, 2, reply.Total)
		}
	})

	t.Run("fail: invalid param", func(t *testing.T) {
		_, err := client.QueryUser(ctx, &pb.QueryUserRequest{Page: -1, Size: -1})
		assert.Error(t, err)
	})
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		_, err := client.CreateUser(ctx, &pb.CreateUserRequest{Username: "newuser", Password: "newpassword"})
		assert.NoError(t, err)
	})

	t.Run("fail: empty field", func(t *testing.T) {
		_, err := client.CreateUser(ctx, &pb.CreateUserRequest{})
		assert.Error(t, err)
	})

	t.Run("fail: db error", func(t *testing.T) {
		_, err := client.CreateUser(ctx, &pb.CreateUserRequest{Username: "error"})
		assert.Error(t, err)
	})
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		_, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{
			Id:       id1.String(),
			Username: "updateduser",
			Password: "updatedpassword"})
		assert.NoError(t, err)
	})

	t.Run("fail: not found", func(t *testing.T) {
		_, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{
			Id:       uuid.NewString(),
			Username: "updateduser",
			Password: "updatedpassword",
		})
		assert.Error(t, err)
	})

	t.Run("fail: invalid uuid", func(t *testing.T) {
		_, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{
			Id:       "xxx",
			Username: "updateduser",
			Password: "updatedpassword",
		})
		assert.Error(t, err)
	})

	t.Run("fail: db error", func(t *testing.T) {
		_, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{
			Id:       uuid.NewString(),
			Username: "error",
			Password: "updatedpassword",
		})
		assert.Error(t, err)
	})
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		_, err := client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: id1.String()})
		assert.NoError(t, err)
	})

	t.Run("fail: not found", func(t *testing.T) {
		_, err := client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: uuid.NewString()})
		assert.Error(t, err)
	})

	t.Run("fail: invalid uuid", func(t *testing.T) {
		_, err := client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: "xxx"})
		assert.Error(t, err)
	})

	t.Run("fail: db error", func(t *testing.T) {
		_, err := client.DeleteUser(ctx, &pb.DeleteUserRequest{})
		assert.Error(t, err)
	})
}
