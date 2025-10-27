package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/alexander-kartavtsev/starship/payment/internal/model"
)

func (s *service) Pay(_ context.Context, req model.PayOrderRequest) (string, error) {
	paymentMethod := req.GetPaymentMethod()

	switch paymentMethod {
	case model.CARD:
		return uuid.NewString(), nil
	case model.CREDIT_CARD:
		return uuid.NewString(), nil
	case model.SBP:
		return uuid.NewString(), nil
	case model.INVESTOR_MONEY:
		return uuid.NewString(), nil
	case model.UNKNOWN:
		return "", model.ErrNotAvailablePaymentMethod
	}
	return "", model.ErrUnexpectedPaymentMethod
}
