package model

import "time"

type User struct {
	Uuid      string
	Info      *UserInfo
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type UserInfo struct {
	Login               string
	Email               string
	NotificationMethods []*NotificationMethod
}

type NotificationMethod struct {
	ProviderName string `json:"provider_name"`
	Target       string `json:"target"`
}

type UserRegistrationInfo struct {
	User     User
	Password string
}
