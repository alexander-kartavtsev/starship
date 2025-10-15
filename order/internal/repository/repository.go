package repository

import (
	"context"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) (string, error)
	Get(ctx context.Context, uuid string) (model.Order, error)
	Update(ctx context.Context, uuid string, updateInfo model.OrderUpdateInfo) error
}

type InventoryRepository interface {
}
