package session

import (
	"context"
	"time"

	"github.com/alexander-kartavtsev/starship/iam/internal/config"
	"github.com/alexander-kartavtsev/starship/iam/internal/converter"
	"github.com/alexander-kartavtsev/starship/iam/internal/model"
)

func (r *repository) Create(ctx context.Context, user *model.User) (*model.Session, error) {
	sessUuid := converter.GenerateUuid()
	now := time.Now()
	session := model.Session{
		SessionUuid: sessUuid,
		UserUuid:    user.Uuid,
		CreatedAt:   now,
		UpdatedAt:   &now,
		ExpiresAt:   &now,
	}

	sessionKey := sessionKeyPrefix + sessUuid

	err := r.cache.HashSet(ctx, sessionKey, converter.SessionToRedisView(session))
	if err != nil {
		return nil, err
	}
	err = r.cache.Expire(ctx, sessionKey, config.AppConfig().Session.SessionTTL())
	if err != nil {
		return nil, err
	}

	return &session, nil
}
