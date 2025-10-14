package order

import (
	"github.com/alexander-kartavtsev/starship/order/internal/repository"
	def "github.com/alexander-kartavtsev/starship/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository
}

func NewService(orderRepository repository.OrderRepository) *service {
	return &service{
		orderRepository: orderRepository,
	}
}
