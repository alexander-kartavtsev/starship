package model

import "errors"

var (
	ErrNotAvailablePaymentMethod error = errors.New("метод оплаты не поддерживается")
	ErrUnexpectedPaymentMethod   error = errors.New("некорректно указан метод оплаты")
)
