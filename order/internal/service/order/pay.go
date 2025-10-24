package order

import (
	"context"
	"log"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (s *service) Pay(ctx context.Context, orderUuid string, payMethod model.PaymentMethod) (string, error) {
	order, err := s.orderRepository.Get(ctx, orderUuid)
	if err != nil {
		return "", err
	}
	log.Println("Оплачиваем...")
	transactionUuid, errPay := s.paymentClient.PayOrder(ctx, model.PayOrderRequest{
		OrderUuid:     orderUuid,
		UserUuid:      order.UserUuid,
		PaymentMethod: payMethod,
	})
	if errPay != nil {
		return "", model.ErrPayment
	}
	log.Println("...готово")
	paidStatus := model.Paid
	log.Println("Обновляем...")

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
		return transactionUuid, err
	}
	log.Println("...готово")
	return transactionUuid, nil
}
