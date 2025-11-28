package session

import (
	"context"

	"github.com/alexander-kartavtsev/starship/iam/internal/model"
)

func (r *repository) Get(ctx context.Context, sessionUuid string) (*model.Whoami, error) {
	return nil, nil
}
