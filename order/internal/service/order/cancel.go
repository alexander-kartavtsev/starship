package order

import (
	"context"

	"github.com/samber/lo"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (s *service) Cansel(ctx context.Context, uuid string) error {
	err := s.orderRepository.Update(
		ctx,
		uuid,
		model.OrderUpdateInfo{
			Status: lo.ToPtr(model.Cancelled),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
