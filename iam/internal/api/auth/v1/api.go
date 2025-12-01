package v1

import (
	"github.com/alexander-kartavtsev/starship/iam/internal/service"
	authV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/auth/v1"
)

type api struct {
	authV1.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewApi(authService service.AuthService) *api {
	return &api{
		authService: authService,
	}
}
