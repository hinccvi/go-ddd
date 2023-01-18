package v1

import (
	"context"

	"github.com/google/uuid"
	"github.com/hinccvi/go-ddd/pkg/log"

	"github.com/hinccvi/go-ddd/internal/user/service"
	"github.com/hinccvi/go-ddd/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type resource struct {
	proto.UnimplementedUserServiceServer
	logger  log.Logger
	service service.Service
}

func RegisterHandlers(service service.Service, logger log.Logger) resource {
	return resource{logger: logger, service: service}
}

func (r resource) GetUser(ctx context.Context, req *proto.GetUserRequest) (reply *proto.GetUserReply, err error) {
	if err = req.Validate(); err != nil {
		return &proto.GetUserReply{}, err
	}

	data := service.GetUserRequest{
		ID: uuid.MustParse(req.Id),
	}
	userEntity, err := r.service.GetUser(ctx, data)
	if err != nil {
		return &proto.GetUserReply{}, err
	}

	reply.User = &proto.User{
		Username:  userEntity.Username,
		CreatedAt: timestamppb.New(userEntity.CreatedAt),
	}

	return reply, nil
}

func (r resource) QueryUser(ctx context.Context, req *proto.QueryUserRequest) (reply *proto.QueryUserReply, err error) {
	if err = req.Validate(); err != nil {
		return &proto.QueryUserReply{}, err
	}

	data := service.QueryUserRequest{
		Page: int(req.Page),
		Size: int(req.Size),
	}
	users, total, err := r.service.QueryUser(ctx, data)
	if err != nil {
		return &proto.QueryUserReply{}, err
	}

	var pbUsers []*proto.User
	for _, user := range users {
		pbUsers = append(pbUsers, &proto.User{
			Username:  user.Username,
			CreatedAt: timestamppb.New(user.CreatedAt),
		})
	}

	reply.Users = pbUsers
	reply.Total = total

	return reply, nil
}

func (r resource) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data := service.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
	}
	if err := r.service.CreateUser(ctx, data); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r resource) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data := service.UpdateUserRequest{
		ID:       uuid.MustParse(req.Id),
		Username: req.Username,
		Password: req.Password,
	}
	if err := r.service.UpdateUser(ctx, data); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r resource) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	data := service.DeleteUserRequest{
		ID: uuid.MustParse(req.Id),
	}
	if err := r.service.DeleteUser(ctx, data); err != nil {
		return nil, err
	}

	return nil, nil
}
