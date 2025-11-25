package app

import (
	"context"

	apiV1 "github.com/alexander-kartavtsev/starship/payment/internal/api/payment/v1"
	"github.com/alexander-kartavtsev/starship/payment/internal/service"
	"github.com/alexander-kartavtsev/starship/payment/internal/service/payment"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
	paymentV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentApi     paymentV1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentApi(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentApi == nil {
		d.paymentApi = apiV1.NewApi(d.PaymentService(ctx))
		logger.Info(ctx, "Инициализация Api")
	}
	return d.paymentApi
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = payment.NewService()
		logger.Info(ctx, "Инициализация Service")
	}
	return d.paymentService
}
