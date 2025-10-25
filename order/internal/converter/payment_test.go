package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	paymentV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/payment/v1"
)

func TestPayOrderRequestToProto(t *testing.T) {
	tests := []struct {
		name                  string
		paymentMethod         model.PaymentMethod
		expectedPaymentMethod paymentV1.PaymentMethod
	}{
		{
			name:                  "payment_method_unknown",
			paymentMethod:         model.Unknown,
			expectedPaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED,
		},
		{
			name:                  "payment_method_card",
			paymentMethod:         model.Card,
			expectedPaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		},
		{
			name:                  "payment_method_sbp",
			paymentMethod:         model.Sbp,
			expectedPaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
		},
		{
			name:                  "payment_method_credit_card",
			paymentMethod:         model.CreditCard,
			expectedPaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		},
		{
			name:                  "payment_method_investor_money",
			paymentMethod:         model.InvestorMoney,
			expectedPaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
		},
	}

	orderUuid := "any_order_uuid"
	userUuid := "any_user_uuid"

	for _, test := range tests {
		req := model.PayOrderRequest{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: test.paymentMethod,
		}

		expectedRes := paymentV1.PayOrderRequest{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: test.expectedPaymentMethod,
		}
		res := PayOrderRequestToProto(req)

		assert.Equal(t, &expectedRes, res)
	}
}
