package part

import (
	"context"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	"github.com/alexander-kartavtsev/starship/inventory/internal/repository/converter"
)

func (r *repository) List(_ context.Context) (map[string]model.Part, error) {
	list := r.data
	if len(list) == 0 {
		return map[string]model.Part{}, model.ErrPartListEmpty
	}
	return converter.PartsToModel(list), nil
}
