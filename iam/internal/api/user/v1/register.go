package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/alexander-kartavtsev/starship/iam/internal/converter"
	userV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	userUuid, err := a.userService.Register(ctx, converter.UserInfoToModel(req.GetInfo().GetInfo()))
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	return &userV1.RegisterResponse{
		UserUuid: userUuid,
	}, nil
}
