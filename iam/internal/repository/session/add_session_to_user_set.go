package session

import (
	"context"

	"github.com/alexander-kartavtsev/starship/iam/internal/config"
	"github.com/alexander-kartavtsev/starship/iam/internal/model"
)

func (r *repository) AddSessionToUserSet(ctx context.Context, session *model.Session) error {
	key := sessionSetKeyPrefix + session.UserUuid
	err := r.cache.Set(ctx, key, session.SessionUuid)
	if err != nil {
		return err
	}
	err = r.cache.Expire(ctx, key, config.AppConfig().Session.SessionTTL())
	if err != nil {
		return err
	}
	return nil
}
