package user

import (
	"context"

	"github.com/alexander-kartavtsev/starship/iam/internal/model"
)

func (s *service) Get(ctx context.Context, userUuid string) (*model.User, error) {
	user, err := s.userRepository.Get(ctx, userUuid)
	if err != nil {
		return nil, err
	}
	return user, nil
}
