package v1

import (
	"github.com/alexander-kartavtsev/starship/payment/internal/service"
	paymentV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

func NewApi(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
