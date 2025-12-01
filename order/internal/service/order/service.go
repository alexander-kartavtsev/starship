package order

import (
	"github.com/alexander-kartavtsev/starship/order/internal/client/grpc"
	"github.com/alexander-kartavtsev/starship/order/internal/repository"
	def "github.com/alexander-kartavtsev/starship/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository      repository.OrderRepository
	inventoryClient      grpc.InventoryClient
	paymentClient        grpc.PaymentClient
	orderProducerService def.OrderProducerService
	orderConsumerService def.ConsumerService
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
	orderProducerService def.OrderProducerService,
	orderConsumerService def.ConsumerService,
) *service {
	return &service{
		orderRepository:      orderRepository,
		inventoryClient:      inventoryClient,
		paymentClient:        paymentClient,
		orderProducerService: orderProducerService,
		orderConsumerService: orderConsumerService,
	}
}
