package session

import (
	"context"

	"github.com/alexander-kartavtsev/starship/iam/internal/model"
)

func (r *repository) Create(ctx context.Context, loginData *model.LoginData) (*model.Session, error) {
	return nil, nil
}
