package v1

import (
	"context"

	"github.com/hinccvi/go-ddd/pkg/log"
	"github.com/hinccvi/go-ddd/proto/pb"
	"google.golang.org/grpc/metadata"

	"github.com/hinccvi/go-ddd/internal/auth/service"
	"github.com/hinccvi/go-ddd/internal/errors"
)

type resource struct {
	pb.UnimplementedAuthServiceServer
	logger  log.Logger
	service service.Service
}

func RegisterHandlers(service service.Service, logger log.Logger) resource {
	return resource{logger: logger, service: service}
}

func (r resource) Login(ctx context.Context, req *pb.LoginRequest) (reply *pb.LoginReply, err error) {
	if req.Username == "" || req.Password == "" {
		err = errors.EmptyField.E()
		return
	}

	data := service.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	accessToken, refreshToken, err := r.service.Login(ctx, data)
	if err != nil {
		return
	}

	reply = &pb.LoginReply{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return
}

func (r resource) Refresh(ctx context.Context, req *pb.RefreshRequest) (reply *pb.RefreshReply, err error) {
	if req.RefreshToken == "" {
		err = errors.EmptyField.E()
		return
	}

	var values []string
	var token string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values = md.Get("Authorization")
	}

	if len(values) > 0 {
		token = values[0]
	}

	data := service.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
		AccessToken:  token,
	}
	accessToken, err := r.service.Refresh(ctx, data)
	if err != nil {
		return
	}

	reply = &pb.RefreshReply{
		AccessToken: accessToken,
	}

	return
}
