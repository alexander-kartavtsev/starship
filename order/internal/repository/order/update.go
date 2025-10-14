package order

import (
	"context"
	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (r *repository) Update(_ context.Context, uuid string, info model.OrderUpdateInfo) error {
	return nil
}
