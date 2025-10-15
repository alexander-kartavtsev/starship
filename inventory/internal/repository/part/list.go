package part

import (
	"context"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	"github.com/alexander-kartavtsev/starship/inventory/internal/repository/converter"
)

func (r *repository) List(_ context.Context, filter model.PartsFilter) (map[string]*model.Part, error) {
	repoFilter := converter.PartsFilterToRepoModel(filter)
	list, err := GetPartsByFilter(r.data, &repoFilter)
	if err != nil {
		return map[string]*model.Part{}, err
	}
	return converter.PartsToModel(list), nil
}
