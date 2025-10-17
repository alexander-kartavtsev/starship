package part

import (
	"context"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	repository "github.com/alexander-kartavtsev/starship/inventory/internal/repository_stub"
)

type service struct {
	inventoryRepository repository.InventoryRepository
}

func NewService(inventoryRepository repository.InventoryRepository) *service {
	return &service{
		inventoryRepository: inventoryRepository,
	}
}

func (s *service) Get(ctx context.Context, partUuid string) (model.Part, error) {
	part, err := s.inventoryRepository.Get(ctx, partUuid)
	if err != nil {
		return model.Part{}, err
	}
	return part, nil
}
