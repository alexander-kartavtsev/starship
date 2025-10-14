package order

import (
	"context"
	"github.com/alexander-kartavtsev/starship/order/internal/model"
	repoConverter "github.com/alexander-kartavtsev/starship/order/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, uuid string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoOrder, ok := r.data[uuid]
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	return repoConverter.OrderToModel(repoOrder), nil
}
