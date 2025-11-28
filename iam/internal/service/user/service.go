package user

import (
	"github.com/alexander-kartavtsev/starship/iam/internal/repository"
	def "github.com/alexander-kartavtsev/starship/iam/internal/service"
)

var _ def.UserService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}
