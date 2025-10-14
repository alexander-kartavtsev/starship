package order

import (
	"context"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (s *service) Create(ctx context.Context, info model.OrderInfo) (*model.OrderCreateRes, error) {
	var totalPrice float64 = 12500

	order := model.Order{
		UserUuid:   info.GetUserUuid(),
		PartUuids:  info.GetPartUuids(),
		TotalPrice: totalPrice,
	}

	orderUuid, err := s.orderRepository.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return &model.OrderCreateRes{
		OrderUuid:  orderUuid,
		TotalPrice: order.TotalPrice,
	}, nil
}
