package v1

import (
	"context"

	"github.com/alexander-kartavtsev/starship/iam/internal/converter"
	"github.com/alexander-kartavtsev/starship/iam/internal/model"
	userV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	userUuid, err := a.userService.Register(
		ctx,
		converter.UserInfoToModel(req.GetInfo().GetInfo()),
		req.GetInfo().GetPassword(),
	)
	if err != nil {
		return nil, model.ErrUserServUnavailable
	}

	return &userV1.RegisterResponse{
		UserUuid: userUuid,
	}, nil
}
