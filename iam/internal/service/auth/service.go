package auth

import (
	"github.com/alexander-kartavtsev/starship/iam/internal/repository"
	def "github.com/alexander-kartavtsev/starship/iam/internal/service"
)

var _ def.AuthService = (*service)(nil)

type service struct {
	sessionRepository repository.SessionRepository
	userRepository    repository.UserRepository
}

func NewService(sessionRepository repository.SessionRepository, userRepository repository.UserRepository) *service {
	return &service{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
	}
}
