package grpc

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/hinccvi/go-ddd/proto/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
)

const bufSize = 1024 * 1024

var (
	lis    *bufconn.Listener
	client pb.HealthcheckServiceClient
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
	client = pb.NewHealthcheckServiceClient(conn)

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
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	pb.RegisterHealthcheckServiceServer(
		s,
		RegisterHandlers("1.0"),
	)

	return s
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGetVersion(t *testing.T) {
	ctx := context.Background()

	reply, err := client.GetVersion(ctx, &emptypb.Empty{})
	if assert.NoError(t, err) {
		assert.Equal(t, "1.0", reply.Version)
	}
}
