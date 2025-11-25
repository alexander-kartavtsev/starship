package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	grpcMocks "github.com/alexander-kartavtsev/starship/order/internal/client/grpc/mocks"
	"github.com/alexander-kartavtsev/starship/order/internal/repository/mocks"
	orderServiceMocks "github.com/alexander-kartavtsev/starship/order/internal/service/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	orderRepository      *mocks.OrderRepository
	inventoryClient      *grpcMocks.InventoryClient
	paymentClient        *grpcMocks.PaymentClient
	orderProducerService *orderServiceMocks.OrderProducerService

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.inventoryClient = grpcMocks.NewInventoryClient(s.T())
	s.paymentClient = grpcMocks.NewPaymentClient(s.T())
	s.orderProducerService = orderServiceMocks.NewOrderProducerService(s.T())

	s.service = NewService(
		s.orderRepository,
		s.inventoryClient,
		s.paymentClient,
		s.orderProducerService,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
