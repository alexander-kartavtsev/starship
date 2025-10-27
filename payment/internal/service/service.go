package service

import (
	"context"

	"github.com/alexander-kartavtsev/starship/payment/internal/model"
)

type PaymentService interface {
	Pay(context.Context, model.PayOrderRequest) (string, error)
}
