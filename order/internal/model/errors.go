package model

import "errors"

var (
	ErrOrderNotFound     = errors.New("ошибка: заказ не найден")
	ErrCancelPaidOrder   = errors.New("ошибка: нельзя отменить оплаченный заказ")
	ErrPartsNotAvailable = errors.New("ошибка: запчасти недоступны")
	ErrPayment           = errors.New("ошибка: не удалось оплатить")
)
