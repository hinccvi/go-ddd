package v1

import (
	"context"

	"github.com/hinccvi/go-ddd/pkg/log"

	"github.com/hinccvi/go-ddd/internal/auth/service"
	"github.com/hinccvi/go-ddd/proto"
)

type resource struct {
	proto.UnimplementedAuthServiceServer
	logger  log.Logger
	service service.Service
}

func RegisterHandlers(service service.Service, logger log.Logger) resource {
	return resource{logger: logger, service: service}
}

func (r resource) Login(ctx context.Context, req *proto.LoginRequest) (reply *proto.LoginReply, err error) {
	data := service.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	accessToken, refreshToken, err := r.service.Login(ctx, data)
	if err != nil {
		return &proto.LoginReply{}, err
	}

	reply.AccessToken = accessToken
	reply.RefreshToken = refreshToken

	return reply, nil
}

func (r resource) Refresh(ctx context.Context, req *proto.RefreshRequest) (reply *proto.RefreshReply, err error) {
	data := service.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}
	accessToken, err := r.service.Refresh(ctx, data)
	if err != nil {
		return reply, err
	}

	reply.AccessToken = accessToken

	return reply, nil
}
