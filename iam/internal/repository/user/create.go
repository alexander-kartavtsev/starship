package user

import (
	"context"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"

	"github.com/alexander-kartavtsev/starship/iam/internal/converter"
	"github.com/alexander-kartavtsev/starship/iam/internal/model"
)

func (r *repository) Create(ctx context.Context, info *model.UserInfo, password string) (string, error) {
	if len(info.Login) == 0 {
		return "", errors.New("логин - обязательное поле")
	}

	user, err := r.GetBy(ctx, "login", info.Login)
	if err != nil {
		if errors.Is(err, model.ErrGetUser) {
			return "", err
		}
	}

	if user != nil {
		log.Printf("пользователь с таким логином уже есть в системе")
		return "", model.ErrUserAlreadyExists
	}

	if len(info.Email) > 0 {
		user, err = r.GetBy(ctx, "email", info.Email)
		if err != nil {
			if errors.Is(err, model.ErrGetUser) {
				return "", err
			}
		}
		if user != nil {
			log.Printf("пользователь с таким email уже есть в системе")
			return "", model.ErrUserAlreadyExists
		}
	}

	userUuid := converter.GenerateUuid()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Printf("Не удалось сгенерировать хэш пароля")
		return "", model.ErrCreateUser
	}

	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("user_uuid", "login", "password", "email", "notifications").
		Values(userUuid, info.Login, string(passwordHash), info.Email, info.NotificationMethods).
		Suffix("RETURNING user_uuid")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("Ошибка build query: %v\n", err)
		return "", model.ErrCreateUser
	}

	var userUuidDb string
	err = r.poolDb.QueryRow(ctx, query, args...).Scan(&userUuidDb)
	if err != nil {
		log.Printf("Ошибка insert в таблицу users: %v\n", err)
		return "", model.ErrCreateUser
	}

	return userUuidDb, nil
}
