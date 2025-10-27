package repository

import (
	"context"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
)

//go:generate ./../../../bin/mockery --case=underscore --all

type InventoryRepository interface {
	Get(context.Context, string) (model.Part, error)
	List(context.Context) (map[string]model.Part, error)
}
