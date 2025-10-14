package service

import (
	"context"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

type OrderService interface {
	Cansel(ctx context.Context, uuid string) error
	Create(ctx context.Context, info model.OrderInfo) (*model.OrderCreateRes, error)
	Get(ctx context.Context, uuid string) (model.Order, error)
	Pay(ctx context.Context, uuid string, payMethod model.PaymentMethod) (string, error)
}
