package part

import (
	"context"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
)

func (s *service) List(ctx context.Context, filter model.PartsFilter) (map[string]model.Part, error) {
	parts, errRep := s.inventoryRepository.List(ctx)
	if errRep != nil {
		return map[string]model.Part{}, errRep
	}
	list, errFound := GetPartsByFilter(parts, &filter)
	if errFound != nil {
		return map[string]model.Part{}, errFound
	}

	return list, nil
}
