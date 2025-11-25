package order

import (
	"context"

	"github.com/samber/lo"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, orderUuid string) error {
	order, err := s.orderRepository.Get(ctx, orderUuid)
	if err != nil {
		return err
	}

	status := order.Status
	if status == model.Paid {
		return model.ErrCancelPaidOrder
	}

	err = s.orderRepository.Update(
		ctx,
		orderUuid,
		model.OrderUpdateInfo{
			Status: lo.ToPtr(model.Cancelled),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
