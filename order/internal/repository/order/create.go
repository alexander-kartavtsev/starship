package order

import (
	"context"
	"github.com/alexander-kartavtsev/starship/order/internal/model"
	"github.com/alexander-kartavtsev/starship/order/internal/repository/converter"
	repoModel "github.com/alexander-kartavtsev/starship/order/internal/repository/model"
	"github.com/google/uuid"
)

func (r *repository) Create(_ context.Context, info model.OrderInfo) (model.Order, error) {
	newUuid := uuid.NewString()

	r.mu.Lock()
	defer r.mu.Unlock()

	repoOrder := repoModel.Order{
		Uuid:       newUuid,
		TotalPrice: 10000,
	}
	r.data[newUuid] = repoOrder

	return converter.OrderToModel(repoOrder), nil
}
