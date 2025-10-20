package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	repoModel "github.com/alexander-kartavtsev/starship/order/internal/repository/model"
)

func (r *repository) Create(_ context.Context, order model.Order) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	orderUuid := uuid.NewString()

	repoOrder := repoModel.Order{
		OrderUuid:  orderUuid,
		UserUuid:   order.UserUuid,
		PartUuids:  order.PartUuids,
		TotalPrice: order.TotalPrice,
		Status:     repoModel.PendingPayment,
	}
	r.data[orderUuid] = repoOrder

	return orderUuid, nil
}
