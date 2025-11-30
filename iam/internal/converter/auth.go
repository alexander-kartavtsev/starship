package converter

import (
	"time"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/alexander-kartavtsev/starship/iam/internal/model"
	repoModel "github.com/alexander-kartavtsev/starship/iam/internal/repository/model"
	commonV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/common/v1"
)

func SessionToProto(session *model.Session) *commonV1.Session {
	if session == nil {
		return nil
	}
	return &commonV1.Session{
		SessionUuid: session.SessionUuid,
		UserUuid:    session.UserUuid,
		CreatedAt:   timestamppb.New(session.CreatedAt),
		UpdatedAt:   timestamppb.New(*session.UpdatedAt),
		ExpiresAt:   timestamppb.New(*session.ExpiresAt),
	}
}

func SessionToRedisView(session model.Session) repoModel.SessionRedisView {
	var updatedAt *int64
	if session.UpdatedAt != nil {
		updatedAt = lo.ToPtr(session.UpdatedAt.UnixNano())
	}

	var expiresAt *int64
	if session.ExpiresAt != nil {
		expiresAt = lo.ToPtr(session.ExpiresAt.UnixNano())
	}

	return repoModel.SessionRedisView{
		SessionUuid: session.SessionUuid,
		UserUuid:    session.UserUuid,
		CreatedAtNs: session.CreatedAt.UnixNano(),
		UpdatedAtNs: updatedAt,
		ExpiresAtNs: expiresAt,
	}
}

// SessionFromRedisView - конвертер из Redis view в модель домена
func SessionFromRedisView(redisView repoModel.SessionRedisView) model.Session {
	var updatedAt *time.Time
	if redisView.UpdatedAtNs != nil {
		tmp := time.Unix(0, *redisView.UpdatedAtNs)
		updatedAt = &tmp
	}

	var expiresAt *time.Time
	if redisView.ExpiresAtNs != nil {
		tmp := time.Unix(0, *redisView.ExpiresAtNs)
		expiresAt = &tmp
	}

	return model.Session{
		SessionUuid: redisView.SessionUuid,
		UserUuid:    redisView.UserUuid,
		CreatedAt:   time.Unix(0, redisView.CreatedAtNs),
		UpdatedAt:   updatedAt,
		ExpiresAt:   expiresAt,
	}
}
