package v1

import (
	"context"

	"github.com/google/uuid"
	"github.com/hinccvi/go-ddd/pkg/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hinccvi/go-ddd/internal/errors"
	"github.com/hinccvi/go-ddd/internal/user/service"
	"github.com/hinccvi/go-ddd/proto/pb"
)

type resource struct {
	pb.UnimplementedUserServiceServer
	logger  log.Logger
	service service.Service
}

func RegisterHandlers(service service.Service, logger log.Logger) resource {
	return resource{logger: logger, service: service}
}

func (r resource) GetUser(ctx context.Context, req *pb.GetUserRequest) (reply *pb.GetUserReply, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	data := service.GetUserRequest{
		ID: uuid.MustParse(req.Id),
	}
	userEntity, err := r.service.GetUser(ctx, data)
	if err != nil {
		return
	}

	reply = &pb.GetUserReply{
		User: &pb.User{
			Id:        userEntity.ID.String(),
			Username:  userEntity.Username,
			CreatedAt: timestamppb.New(userEntity.CreatedAt),
		},
	}

	return
}

func (r resource) QueryUser(ctx context.Context, req *pb.QueryUserRequest) (reply *pb.QueryUserReply, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	data := service.QueryUserRequest{
		Page: int(req.Page),
		Size: int(req.Size),
	}
	users, total, err := r.service.QueryUser(ctx, data)
	if err != nil {
		return
	}

	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:        user.ID.String(),
			Username:  user.Username,
			CreatedAt: timestamppb.New(user.CreatedAt),
		})
	}

	reply = &pb.QueryUserReply{
		Users: pbUsers,
		Total: total,
	}

	return
}

func (r resource) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (e *emptypb.Empty, err error) {
	if req.Username == "" || req.Password == "" {
		err = errors.EmptyField.E()
		return
	}

	data := service.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
	}
	if err = r.service.CreateUser(ctx, data); err != nil {
		return
	}

	e = &emptypb.Empty{}

	return
}

func (r resource) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (e *emptypb.Empty, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	data := service.UpdateUserRequest{
		ID:       uuid.MustParse(req.Id),
		Username: req.Username,
		Password: req.Password,
	}
	if err = r.service.UpdateUser(ctx, data); err != nil {
		return
	}

	e = &emptypb.Empty{}

	return
}

func (r resource) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (e *emptypb.Empty, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	data := service.DeleteUserRequest{
		ID: uuid.MustParse(req.Id),
	}
	if err = r.service.DeleteUser(ctx, data); err != nil {
		return
	}

	e = &emptypb.Empty{}

	return
}
