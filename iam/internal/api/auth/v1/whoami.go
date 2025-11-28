package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/alexander-kartavtsev/starship/iam/internal/converter"
	authV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, request *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	whoami, err := a.authService.Whoami(ctx, request.GetSessionUuid())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return &authV1.WhoamiResponse{
		Session: converter.SessionToProto(&whoami.Session),
		User:    converter.UserToProto(&whoami.User),
	}, nil
}
