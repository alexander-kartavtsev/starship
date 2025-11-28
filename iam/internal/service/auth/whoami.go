package auth

import (
	"context"

	"github.com/alexander-kartavtsev/starship/iam/internal/model"
)

func (s *service) Whoami(ctx context.Context, sessionUuid string) (*model.Whoami, error) {
	whoami, err := s.sessionRepository.Get(ctx, sessionUuid)
	if err != nil {
		return nil, err
	}
	return whoami, nil
}
