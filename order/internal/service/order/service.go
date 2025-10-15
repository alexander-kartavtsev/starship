package order

import (
	"github.com/alexander-kartavtsev/starship/order/internal/client/grpc"
	"github.com/alexander-kartavtsev/starship/order/internal/repository"
	def "github.com/alexander-kartavtsev/starship/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository

	inventoryClient grpc.InventoryClient
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
	}
}
