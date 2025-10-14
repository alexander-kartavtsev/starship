package order

import (
	"context"
	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (s *service) Pay(ctx context.Context, uuid string, payMethod model.PaymentMethod) (string, error) {
	return "", nil
}
