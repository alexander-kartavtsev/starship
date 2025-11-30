package auth

import (
	"context"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/alexander-kartavtsev/starship/iam/internal/model"
)

func (s *service) Login(ctx context.Context, loginData *model.LoginData) (string, error) {
	user, err := s.userRepository.GetBy(ctx, "login", loginData.Login)
	if err != nil {
		return "", errors.New("login")
		// return "", model.ErrUserNotFound
	}

	log.Printf("введен пароль: %v", loginData.Password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))

	log.Printf("hash введенного: %v, hash пароля из б/д: %v", user.Password, user.Password)

	if err != nil {
		return "", errors.New("password")
		// return "", model.ErrUserNotFound
	}
	session, err := s.sessionRepository.Create(ctx, &user.User)
	if err != nil {
		return "", err
	}

	err = s.sessionRepository.AddSessionToUserSet(ctx, session)
	if err != nil {
		return "", err
	}

	return session.SessionUuid, nil
}
