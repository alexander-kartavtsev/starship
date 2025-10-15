package service

import (
	"context"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
)

type InventoryService interface {
	Get(context.Context, string) (model.Part, error)
	List(context.Context, model.PartsFilter) (map[string]*model.Part, error)
}
