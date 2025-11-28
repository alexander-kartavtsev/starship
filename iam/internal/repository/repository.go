package repository

import (
	"context"

	"github.com/alexander-kartavtsev/starship/iam/internal/model"
)

type SessionRepository interface {
	AddSessionToUserSet(context.Context, *model.Session) error
	Create(context.Context, *model.LoginData) (*model.Session, error)
	Get(context.Context, string) (*model.Whoami, error)
}

type UserRepository interface {
	Create(context.Context, *model.UserInfo) (string, error)
	Get(context.Context, string) (*model.User, error)
}
