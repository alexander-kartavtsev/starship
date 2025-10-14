package order

import "github.com/alexander-kartavtsev/starship/order/internal/repository"
import def "github.com/alexander-kartavtsev/starship/order/internal/service"

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepsitory repository.OrderRepository
}

func NewService(orderRepository repository.OrderRepository) *service {
	return &service{
		orderRepsitory: orderRepository,
	}
}
