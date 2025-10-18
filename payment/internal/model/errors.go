package model

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotAvailablePaymentMethod  error = errors.New("метод оплаты не поддерживается")
	ErrUnexpectedPaymentMethod    error = errors.New("неизвестный метод оплаты")
	ErrNotAvailablePayMethodProto error = status.Errorf(codes.InvalidArgument, "Способ оплаты не поддерживается")
	ErrUnexpectedPaytMethodProto  error = status.Errorf(codes.InvalidArgument, "Неизвестный метод оплаты")
)
