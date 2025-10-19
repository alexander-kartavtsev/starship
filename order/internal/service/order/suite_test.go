package order

import (
	"context"
	grpcMocks "github.com/alexander-kartavtsev/starship/order/internal/client/grpc/mocks"
	"github.com/alexander-kartavtsev/starship/order/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	orderRepository *mocks.OrderRepository
	inventoryClient *grpcMocks.InventoryClient
	paymentClient   *grpcMocks.PaymentClient

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.inventoryClient = grpcMocks.NewInventoryClient(s.T())
	s.paymentClient = grpcMocks.NewPaymentClient(s.T())

	s.service = NewService(
		s.orderRepository,
		s.inventoryClient,
		s.paymentClient,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
