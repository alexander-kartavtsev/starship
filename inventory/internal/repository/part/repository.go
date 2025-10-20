package part

import (
	"sync"

	def "github.com/alexander-kartavtsev/starship/inventory/internal/repository"
	repoModel "github.com/alexander-kartavtsev/starship/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Part
}

func NewRepository() *repository {
	return &repository{
		data: GetAllParts(),
	}
}
