package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/alexander-kartavtsev/starship/iam/internal/model"
	commonV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/common/v1"
)

func SessionToProto(session *model.Session) *commonV1.Session {
	if session == nil {
		return nil
	}
	return &commonV1.Session{
		SessionUuid: session.SessionUuid,
		CreatedAt:   timestamppb.New(*session.CreatedAt),
		UpdatedAt:   timestamppb.New(*session.UpdatedAt),
		ExpiresAt:   timestamppb.New(*session.ExpiresAt),
	}
}
