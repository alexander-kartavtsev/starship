package order

import (
	"context"
	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (s *service) Cansel(ctx context.Context, uuid string) error {
	err := s.orderRepsitory.Update(
		ctx,
		uuid,
		model.OrderUpdateInfo{
			Status: model.Cancelled,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
