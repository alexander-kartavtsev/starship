package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/alexander-kartavtsev/starship/iam/internal/converter"
	userV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/user/v1"
)

func (a *api) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	user, err := a.userService.Get(ctx, req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &userV1.GetUserResponse{
		User: converter.UserToProto(user),
	}, nil
}
