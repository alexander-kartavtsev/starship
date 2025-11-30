package v1

import (
	"context"

	"github.com/alexander-kartavtsev/starship/iam/internal/converter"
	"github.com/alexander-kartavtsev/starship/iam/internal/model"
	userV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/user/v1"
)

func (a *api) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	user, err := a.userService.Get(ctx, req.GetUserUuid())
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	return &userV1.GetUserResponse{
		User: converter.UserToProto(user),
	}, nil
}
