package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hinccvi/go-ddd/proto/pb"
)

type resource struct {
	pb.UnimplementedHealthcheckServiceServer
	version string
}

func RegisterHandlers(version string) resource {
	return resource{version: version}
}

func (r resource) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, nil
}

func (r resource) GetVersion(ctx context.Context, empty *emptypb.Empty) (*pb.GetVersionReply, error) {
	return &pb.GetVersionReply{Version: r.version}, nil
}
