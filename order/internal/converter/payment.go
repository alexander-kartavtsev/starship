package converter

import (
	"github.com/alexander-kartavtsev/starship/order/internal/model"
	paymentV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/payment/v1"
)

func PayOrderRequestToProto(req model.PayOrderRequest) *paymentV1.PayOrderRequest {
	return &paymentV1.PayOrderRequest{
		UserUuid:      req.UserUuid,
		OrderUuid:     req.OrderUuid,
		PaymentMethod: PaymentMethodToProto(req.PaymentMethod),
	}
}

func PaymentMethodToProto(paymentMethod model.PaymentMethod) paymentV1.PaymentMethod {
	var res paymentV1.PaymentMethod
	switch paymentMethod {
	case model.Unknown:
		res = paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED
	case model.Card:
		res = paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.Sbp:
		res = paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.CreditCard:
		res = paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.InvestorMoney:
		res = paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	}
	return res
}
