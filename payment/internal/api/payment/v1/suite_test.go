package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/alexander-kartavtsev/starship/payment/internal/service/mocks"
)

type ApiSuite struct {
	suite.Suite
	ctx            context.Context
	paymentService *mocks.PaymentService
	api            *api
}

func (a *ApiSuite) SetupTest() {
	a.ctx = context.Background()
	a.paymentService = mocks.NewPaymentService(a.T())
	a.api = NewApi(a.paymentService)
}

func (a *ApiSuite) TearDownTest() {
}

func TestApiIntegration(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
