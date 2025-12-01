package user

import (
	"context"

	"github.com/alexander-kartavtsev/starship/iam/internal/model"
)

func (s *service) Register(ctx context.Context, info *model.UserInfo, password string) (string, error) {
	userUuid, err := s.userRepository.Create(ctx, info, password)
	if err != nil {
		return "", err
	}
	return userUuid, nil
}
