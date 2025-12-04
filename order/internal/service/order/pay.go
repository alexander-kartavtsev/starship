package order

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	"github.com/alexander-kartavtsev/starship/platform/pkg/tracing"
)

func (s *service) Pay(ctx context.Context, orderUuid string, payMethod model.PaymentMethod) (string, error) {
	ctx, spanOrderGet := tracing.StartSpan(ctx, "order.repository.Get",
		trace.WithAttributes(
			attribute.String("order_uuid", orderUuid),
		),
	)

	order, err := s.orderRepository.Get(ctx, orderUuid)
	if err != nil {
		spanOrderGet.RecordError(err)
		spanOrderGet.End()
		return "", err
	}
	spanOrderGet.End()

	ctx, spanPayOrder := tracing.StartSpan(ctx, "order.paymentClient.PayOrder",
		trace.WithAttributes(
			attribute.String("order_uuid", orderUuid),
			attribute.String("user_uuid", order.UserUuid),
			attribute.String("payment_method", string(payMethod)),
		),
	)
	log.Println("Оплачиваем...")
	transactionUuid, errPay := s.paymentClient.PayOrder(ctx, model.PayOrderRequest{
		OrderUuid:     orderUuid,
		UserUuid:      order.UserUuid,
		PaymentMethod: payMethod,
	})
	if errPay != nil {
		spanPayOrder.RecordError(err)
		spanPayOrder.End()
		log.Printf("Ошибка paymentClient.PayOrder: %v\n", errPay)
		return "", model.ErrPayment
	}
	log.Println("...готово")
	spanPayOrder.End()

	paidStatus := model.Paid
	log.Println("Обновляем...")

	ctx, spanUpdateOrder := tracing.StartSpan(ctx, "order.repository.Update",
		trace.WithAttributes(
			attribute.String("order_uuid", orderUuid),
			attribute.String("status", string(paidStatus)),
			attribute.String("transaction_uuid", transactionUuid),
			attribute.String("payment_method", string(payMethod)),
		),
	)

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
		spanUpdateOrder.RecordError(err)
		spanUpdateOrder.End()
		return transactionUuid, err
	}
	log.Println("...готово")
	spanUpdateOrder.End()

	ctx, spanProduceToKafka := tracing.StartSpan(ctx, "order.orderProducerService.ProduceOrder",
		trace.WithAttributes(
			attribute.String("order_uuid", orderUuid),
		),
	)
	defer spanProduceToKafka.End()

	err = s.orderProducerService.ProduceOrder(ctx, model.OrderKafkaEvent{
		Uuid:            orderUuid,
		OrderUuid:       order.OrderUuid,
		UserUuid:        order.UserUuid,
		PaymentMethod:   order.PaymentMethod,
		TransactionUuid: transactionUuid,
		Type:            "pay",
	})
	if err != nil {
		spanProduceToKafka.RecordError(err)
		log.Printf("Ошибка orderProducerService.ProduceOrder: %v\n", err)
		return "", err
	}
	log.Printf("Опубликовали: %v\n", s.orderProducerService)

	return transactionUuid, nil
}
