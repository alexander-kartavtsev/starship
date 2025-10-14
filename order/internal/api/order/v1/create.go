package v1

import (
	"context"
	"github.com/alexander-kartavtsev/starship/order/internal/repository/converter"
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	res, err := a.orderService.Create(
		ctx,
		converter.OrderInfoToModel(req.GetUserUUID(), req.GetPartUuids()),
	)
	if err != nil {
		return nil, err
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  res.OrderUuid,
		TotalPrice: res.TotalPrice,
	}, nil
}
