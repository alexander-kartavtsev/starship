package order

import (
	"sync"

	repoModel "github.com/alexander-kartavtsev/starship/order/internal/repository/model"
)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Order
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]repoModel.Order),
	}
}
