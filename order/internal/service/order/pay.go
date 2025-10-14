package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (s *service) Pay(ctx context.Context, orderUuid string, payMethod model.PaymentMethod) (string, error) {
	_, err := s.orderRepository.Get(ctx, orderUuid)
	if err != nil {
		return "", err
	}
	transactionUuid := uuid.NewString()

	paidStatus := model.Paid

	err = s.orderRepository.Update(
		ctx,
		orderUuid,
		model.OrderUpdateInfo{
			Status:          &paidStatus,
			TransactionUuid: &transactionUuid,
			PaymentMethod:   &payMethod,
		},
	)
	if err != nil {
		return "", err
	}

	return transactionUuid, nil
}
