package order

import (
	"errors"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	"github.com/alexander-kartavtsev/starship/platform/pkg/tracing"
)

func (s *ServiceSuite) TestService_Pay_OrderNotFound() {
	orderUuid := "any_order_uuid"

	ctx, spanOrderGet := tracing.StartSpan(s.ctx, "")
	spanOrderGet.End()

	s.orderRepository.
		On("Get", ctx, orderUuid).
		Return(model.Order{}, model.ErrOrderNotFound).
		Once()

	res, err := s.service.Pay(s.ctx, orderUuid, model.Card)
	s.Assert().Empty(res)
	s.Assert().Equal(err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestService_Pay_PaymentError() {
	orderUuid := "any_order_uuid"
	userUuid := "any_user_uuid"
	paymentMethod := model.Card
	order := model.Order{
		OrderUuid:     orderUuid,
		UserUuid:      userUuid,
		PaymentMethod: paymentMethod,
	}
	payOrderReq := model.PayOrderRequest{
		OrderUuid:     orderUuid,
		UserUuid:      userUuid,
		PaymentMethod: paymentMethod,
	}
	testErr := errors.New("testErr")

	ctx, spanOrderGet := tracing.StartSpan(s.ctx, "")
	spanOrderGet.End()

	s.orderRepository.
		On("Get", ctx, orderUuid).
		Return(order, nil).
		Once()

	ctx, spanPayOrder := tracing.StartSpan(ctx, "")
	spanPayOrder.End()

	s.paymentClient.
		On("PayOrder", ctx, payOrderReq).
		Return("", testErr).
		Once()

	res, err := s.service.Pay(s.ctx, orderUuid, paymentMethod)
	s.Assert().Empty(res)
	s.Assert().True(errors.Is(err, model.ErrPayment))
}

func (s *ServiceSuite) TestService_Pay_UpdateError() {
	orderUuid := "any_order_uuid"
	userUuid := "any_user_uuid"
	paymentMethod := model.Card
	order := model.Order{
		OrderUuid:     orderUuid,
		UserUuid:      userUuid,
		PaymentMethod: paymentMethod,
	}
	payOrderReq := model.PayOrderRequest{
		OrderUuid:     orderUuid,
		UserUuid:      userUuid,
		PaymentMethod: paymentMethod,
	}
	transactionUuid := "any_transaction_uuid"
	paidStatus := model.Paid
	orderUpdateInfo := model.OrderUpdateInfo{
		Status:          &paidStatus,
		TransactionUuid: &transactionUuid,
		PaymentMethod:   &paymentMethod,
	}
	testErr := errors.New("testErr")

	ctx, spanOrderGet := tracing.StartSpan(s.ctx, "")
	spanOrderGet.End()

	s.orderRepository.
		On("Get", ctx, orderUuid).
		Return(order, nil).
		Once()

	ctx, spanPayOrder := tracing.StartSpan(ctx, "")
	spanPayOrder.End()

	s.paymentClient.
		On("PayOrder", ctx, payOrderReq).
		Return(transactionUuid, nil).
		Once()

	ctx, spanUpdateOrder := tracing.StartSpan(ctx, "")
	spanUpdateOrder.End()

	s.orderRepository.
		On("Update", ctx, orderUuid, orderUpdateInfo).
		Return(testErr).
		Once()

	res, err := s.service.Pay(s.ctx, orderUuid, paymentMethod)
	s.Assert().Equal(res, transactionUuid)
	s.Assert().True(errors.Is(err, testErr))
}

func (s *ServiceSuite) TestService_Pay_Ok_Kafka_Fail() {
	orderUuid := "any_order_uuid"
	userUuid := "any_user_uuid"
	paymentMethod := model.Card
	order := model.Order{
		OrderUuid:     orderUuid,
		UserUuid:      userUuid,
		PaymentMethod: paymentMethod,
	}
	payOrderReq := model.PayOrderRequest{
		OrderUuid:     orderUuid,
		UserUuid:      userUuid,
		PaymentMethod: paymentMethod,
	}
	transactionUuid := "any_transaction_uuid"
	paidStatus := model.Paid
	orderUpdateInfo := model.OrderUpdateInfo{
		Status:          &paidStatus,
		TransactionUuid: &transactionUuid,
		PaymentMethod:   &paymentMethod,
	}
	orderKafkaEvent := model.OrderKafkaEvent{
		Uuid:            orderUuid,
		OrderUuid:       orderUuid,
		UserUuid:        userUuid,
		PaymentMethod:   paymentMethod,
		TransactionUuid: transactionUuid,
		Type:            "pay",
	}

	ctx, spanOrderGet := tracing.StartSpan(s.ctx, "")
	spanOrderGet.End()

	s.orderRepository.
		On("Get", ctx, orderUuid).
		Return(order, nil).
		Once()

	ctx, spanPayOrder := tracing.StartSpan(ctx, "")
	spanPayOrder.End()

	s.paymentClient.
		On("PayOrder", ctx, payOrderReq).
		Return(transactionUuid, nil).
		Once()

	ctx, spanUpdateOrder := tracing.StartSpan(ctx, "")
	spanUpdateOrder.End()

	s.orderRepository.
		On("Update", ctx, orderUuid, orderUpdateInfo).
		Return(nil).
		Once()

	ctx, spanProduceToKafka := tracing.StartSpan(ctx, "")
	spanProduceToKafka.End()

	s.orderProducerService.
		On("ProduceOrder", ctx, orderKafkaEvent).
		Return(errors.New("any error")).
		Once()

	res, err := s.service.Pay(s.ctx, orderUuid, paymentMethod)
	s.Assert().Equal(res, "")
	s.Assert().Equal(errors.New("any error"), err)
}

func (s *ServiceSuite) TestService_Pay_Ok() {
	orderUuid := "any_order_uuid"
	userUuid := "any_user_uuid"
	paymentMethod := model.Card
	order := model.Order{
		OrderUuid:     orderUuid,
		UserUuid:      userUuid,
		PaymentMethod: paymentMethod,
	}
	payOrderReq := model.PayOrderRequest{
		OrderUuid:     orderUuid,
		UserUuid:      userUuid,
		PaymentMethod: paymentMethod,
	}
	transactionUuid := "any_transaction_uuid"
	paidStatus := model.Paid
	orderUpdateInfo := model.OrderUpdateInfo{
		Status:          &paidStatus,
		TransactionUuid: &transactionUuid,
		PaymentMethod:   &paymentMethod,
	}
	orderKafkaEvent := model.OrderKafkaEvent{
		Uuid:            orderUuid,
		OrderUuid:       orderUuid,
		UserUuid:        userUuid,
		PaymentMethod:   paymentMethod,
		TransactionUuid: transactionUuid,
		Type:            "pay",
	}

	ctx, spanOrderGet := tracing.StartSpan(s.ctx, "")
	spanOrderGet.End()

	s.orderRepository.
		On("Get", ctx, orderUuid).
		Return(order, nil).
		Once()

	ctx, spanPayOrder := tracing.StartSpan(ctx, "")
	spanPayOrder.End()

	s.paymentClient.
		On("PayOrder", ctx, payOrderReq).
		Return(transactionUuid, nil).
		Once()

	ctx, spanUpdateOrder := tracing.StartSpan(ctx, "")
	spanUpdateOrder.End()

	s.orderRepository.
		On("Update", ctx, orderUuid, orderUpdateInfo).
		Return(nil).
		Once()

	ctx, spanProduceToKafka := tracing.StartSpan(ctx, "")
	spanProduceToKafka.End()

	s.orderProducerService.
		On("ProduceOrder", ctx, orderKafkaEvent).
		Return(nil).
		Once()

	res, err := s.service.Pay(s.ctx, orderUuid, paymentMethod)
	s.Assert().Equal(res, transactionUuid)
	s.Assert().Nil(err)
}
