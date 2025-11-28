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
	ProviderName string
	Target       string
}
