package part

import (
	"context"
	"errors"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	"github.com/alexander-kartavtsev/starship/inventory/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.data[uuid]
	if !ok {
		return model.Part{}, errors.New("not found")
	}
	return converter.PartToModel(*part), nil
}
