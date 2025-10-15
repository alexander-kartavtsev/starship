package grpc

import (
	"context"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) (map[string]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUuid, userUuid, paymentMethod string) (string, error)
}
