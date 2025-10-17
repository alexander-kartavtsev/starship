package part

import (
	"context"
	"errors"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	"github.com/alexander-kartavtsev/starship/inventory/internal/repository/converter"
	repoModel "github.com/alexander-kartavtsev/starship/inventory/internal/repository/model"
	repoPart "github.com/alexander-kartavtsev/starship/inventory/internal/repository/part"
	def "github.com/alexander-kartavtsev/starship/inventory/internal/repository_stub"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	data map[string]*repoModel.Part
}

func NewRepository() *repository {
	return &repository{
		data: map[string]*repoModel.Part{
			"part_uuid_1": {
				Uuid:       "part_uuid_1",
				Dimensions: &repoModel.Dimensions{},
			},
		},
	}
}

func (r *repository) Get(ctx context.Context, uuid string) (model.Part, error) {
	part, ok := r.data[uuid]
	if !ok {
		return model.Part{}, errors.New("not found")
	}
	return converter.PartToModel(*part), nil
}

func (r *repository) List(_ context.Context, filter model.PartsFilter) (map[string]*model.Part, error) {
	repoFilter := converter.PartsFilterToRepoModel(filter)
	list, err := repoPart.GetPartsByFilter(r.data, &repoFilter)
	if err != nil {
		return map[string]*model.Part{}, err
	}
	return converter.PartsToModel(list), nil
}
