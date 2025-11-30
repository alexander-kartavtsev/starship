package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/alexander-kartavtsev/starship/iam/internal/model"
)

func (r *repository) Get(ctx context.Context, userUuid string) (*model.User, error) {
	user, err := getUserBy(ctx, r, "user_uuid", userUuid)
	if err != nil {
		return nil, err
	}
	return &user.User, nil
}

func (r *repository) GetBy(ctx context.Context, field, value string) (*model.UserRegistrationInfo, error) {
	userInfo, err := getUserBy(ctx, r, field, value)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func getUserBy(ctx context.Context, r *repository, field, value string) (*model.UserRegistrationInfo, error) {
	builderQuery := sq.Select(
		"user_uuid",
		"login",
		"password",
		"email",
		"notifications",
		"created_at",
		"updated_at",
	).
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{field: value}).
		Limit(1)

	query, args, err := builderQuery.ToSql()
	if err != nil {
		log.Printf("Ошибка генерации запроса SELECT (Get): %v\n", err)
		return nil, model.ErrGetUser
	}

	var userUuid, login, password string
	var email sql.NullString
	var notification json.RawMessage
	var createdAt, updatedAt *time.Time

	err = r.poolDb.
		QueryRow(ctx, query, args...).
		Scan(
			&userUuid,
			&login,
			&password,
			&email,
			&notification,
			&createdAt,
			&updatedAt,
		)
	if err != nil {
		log.Printf("Ошибка получения данных из таблицы users: %v\n", err)
		if err.Error() == "no rows in result set" {
			return nil, model.ErrUserNotFound
		}
		return nil, model.ErrGetUser
	}

	var restoredNotifications []*model.NotificationMethod

	if notification != nil {
		err = json.Unmarshal(notification, &restoredNotifications)
		if err != nil {
			log.Printf("Ошибка получения данных из таблицы users: %v\n", err)
			return nil, model.ErrGetUser
		}
	}

	user := model.User{
		Uuid: userUuid,
		Info: &model.UserInfo{
			Login:               login,
			Email:               email.String,
			NotificationMethods: restoredNotifications,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return &model.UserRegistrationInfo{
		User:     user,
		Password: password,
	}, nil
}
