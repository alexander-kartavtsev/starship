package model

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserNotFound        = status.Errorf(codes.NotFound, "ошибка: пользователь не найден")
	ErrGetUser             = status.Errorf(codes.Unknown, "ошибка: не удалось получить пользователя")
	ErrUserServUnavailable = status.Errorf(codes.Unavailable, "ошибка: сервис недоступен")
	ErrUserAlreadyExists   = status.Errorf(codes.AlreadyExists, "такой пользователь уже есть")
	ErrCreateUser          = status.Errorf(codes.Unknown, "ошибка: не удалось создать пользователя")

	ErrSessionNotFound = status.Errorf(codes.NotFound, "ошибка: сессия не найдена")
	ErrBadRequest      = status.Errorf(codes.InvalidArgument, "некорректные данные")
)
