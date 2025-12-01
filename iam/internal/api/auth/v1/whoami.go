package v1

import (
	"context"

	"github.com/alexander-kartavtsev/starship/iam/internal/converter"
	"github.com/alexander-kartavtsev/starship/iam/internal/model"
	authV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, request *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	whoami, err := a.authService.Whoami(ctx, request.GetSessionUuid())
	if err != nil {
		return nil, model.ErrSessionNotFound
	}

	sess := converter.SessionToProto(&whoami.Session)
	sess.Reset()

	return &authV1.WhoamiResponse{
		Session: converter.SessionToProto(&whoami.Session),
		User:    converter.UserToProto(&whoami.User),
	}, nil
}
