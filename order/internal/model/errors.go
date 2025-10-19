package model

import "errors"

var (
	ErrOrderNotFound     = errors.New("заказ не найден")
	ErrCancelPaidOrder   = errors.New("нельзя отменить оплаченный заказ")
	ErrPartsNotAvailable = errors.New("запчасти недоступны")
	ErrPayment           = errors.New("ошибка оплаты")
)
