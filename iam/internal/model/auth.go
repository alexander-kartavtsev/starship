package model

import "time"

type LoginData struct {
	Login    string
	Password string
}

type Whoami struct {
	Session Session
	User    User
}

type Session struct {
	SessionUuid string
	UserUuid    string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	ExpiresAt   *time.Time
}
