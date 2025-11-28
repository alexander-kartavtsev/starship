package auth

import (
	"context"

	"github.com/alexander-kartavtsev/starship/iam/internal/model"
)

func (s *service) Login(ctx context.Context, loginData *model.LoginData) (string, error) {
	session, err := s.sessionRepository.Create(ctx, loginData)
	if err != nil {
		return "", err
	}

	err = s.sessionRepository.AddSessionToUserSet(ctx, session)
	if err != nil {
		return "", err
	}

	return session.SessionUuid, nil
}
