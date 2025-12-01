package model

// SessionRedisView - модель для хранения в Redis hash map
type SessionRedisView struct {
	SessionUuid string `redis:"session_uuid"`
	UserUuid    string `redis:"user_uuid"`
	CreatedAtNs int64  `redis:"created_at"`
	UpdatedAtNs *int64 `redis:"updated_at,omitempty"`
	ExpiresAtNs *int64 `redis:"expires_at,omitempty"`
}
