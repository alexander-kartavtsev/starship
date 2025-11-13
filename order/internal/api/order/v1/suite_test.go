package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/alexander-kartavtsev/starship/order/internal/service/mocks"
)

type ApiSuite struct {
	suite.Suite
	ctx          context.Context
	orderService *mocks.OrderService
	api          *api
}

func (a *ApiSuite) SetupTest() {
	a.ctx = context.Background()
	a.orderService = mocks.NewOrderService(a.T())
	a.api = NewApi(a.orderService)
}

func (a *ApiSuite) TearDownTest() {
}

func TestApiIntegration(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
