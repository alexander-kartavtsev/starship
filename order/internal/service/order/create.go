package order

import (
	"context"
	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (s *service) Create(ctx context.Context, info model.OrderInfo) (*model.OrderCreateRes, error) {
	order, err := s.orderRepsitory.Create(ctx, info)

	if err != nil {
		return nil, err
	}

	return &model.OrderCreateRes{
		OrderUuid:  order.OrderUuid,
		TotalPrice: order.TotalPrice,
	}, nil
}
