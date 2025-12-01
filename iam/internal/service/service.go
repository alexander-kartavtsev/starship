package service

import (
	"context"

	"github.com/alexander-kartavtsev/starship/iam/internal/model"
)

type AuthService interface {
	Login(context.Context, *model.LoginData) (string, error)
	Whoami(context.Context, string) (*model.Whoami, error)
}

type UserService interface {
	Get(context.Context, string) (*model.User, error)
	Register(context.Context, *model.UserInfo, string) (string, error)
}
