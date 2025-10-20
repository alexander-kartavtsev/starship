package order

import (
	"context"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.Order, error) {
	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}
