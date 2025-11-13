package order

import (
	"context"

	"github.com/samber/lo"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, uuid string) error {
	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		return err
	}

	status := order.Status
	if status == model.Paid {
		return model.ErrCancelPaidOrder
	}

	err = s.orderRepository.Update(
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
