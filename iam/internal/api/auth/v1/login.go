package v1

import (
	"context"

	"github.com/alexander-kartavtsev/starship/iam/internal/model"
	authV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, request *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	sessionUuid, err := a.authService.Login(ctx, &model.LoginData{
		Login:    request.GetLogin(),
		Password: request.GetPassword(),
	})
	if err != nil {
		return nil, model.ErrBadRequest
	}
	return &authV1.LoginResponse{
		SessionUuid: sessionUuid,
	}, nil
}
