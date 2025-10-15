package part

import (
	"context"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
)

func (s *service) List(ctx context.Context, filter model.PartsFilter) (map[string]*model.Part, error) {
	parts, err := s.inventoryRepository.List(ctx, filter)
	if err != nil {
		return map[string]*model.Part{}, err
	}
	return parts, nil
}
