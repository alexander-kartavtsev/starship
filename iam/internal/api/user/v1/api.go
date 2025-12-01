package v1

import (
	"github.com/alexander-kartavtsev/starship/iam/internal/service"
	userV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/user/v1"
)

type api struct {
	userV1.UnimplementedUserServiceServer
	userService service.UserService
}

func NewApi(userService service.UserService) *api {
	return &api{
		userService: userService,
	}
}
