package session

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	redigo "github.com/gomodule/redigo/redis"

	"github.com/alexander-kartavtsev/starship/iam/internal/converter"
	"github.com/alexander-kartavtsev/starship/iam/internal/model"
	repoModel "github.com/alexander-kartavtsev/starship/iam/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, sessionUuid string) (*model.Whoami, error) {
	sess, err := r.cache.HGetAll(ctx, sessionKeyPrefix+sessionUuid)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return nil, model.ErrSessionNotFound
		}
		return nil, err
	}

	if len(sess) == 0 {
		return nil, model.ErrSessionNotFound
	}

	var sessRedis repoModel.SessionRedisView
	err = redigo.ScanStruct(sess, &sessRedis)
	if err != nil {
		return nil, err
	}

	session := converter.SessionFromRedisView(sessRedis)

	user, err := getUser(ctx, r, session.UserUuid)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	whoami := model.Whoami{
		Session: session,
		User:    user.User,
	}
	return &whoami, nil
}

func getUser(ctx context.Context, r *repository, uuid string) (*model.UserRegistrationInfo, error) {
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
		Where(sq.Eq{"user_uuid": uuid}).
		Limit(1)

	query, args, err := builderQuery.ToSql()
	if err != nil {
		log.Printf("Ошибка генерации запроса SELECT (Get): %v\n", err)
		return nil, model.ErrUserNotFound
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
			return nil, err
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
