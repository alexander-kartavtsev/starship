package v1

import (
	"github.com/alexander-kartavtsev/starship/payment/internal/model"
	"github.com/alexander-kartavtsev/starship/payment/internal/service/payment"
	paymentV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/payment/v1"
)

func (a *ApiSuite) TestApi_PayOrder() {
	const newPaymentMethod paymentV1.PaymentMethod = 5

	api := NewApi(payment.NewService())
	tests := []struct {
		name      string
		payMethod paymentV1.PaymentMethod
		err       error
		expected  int
	}{
		{
			name:      "payment_method_card",
			payMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
			err:       nil,
			expected:  36,
		},
		{
			name:      "payment_method_sbp",
			payMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_SBP,
			err:       nil,
			expected:  36,
		},
		{
			name:      "payment_method_credit_card",
			payMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
			err:       nil,
			expected:  36,
		},
		{
			name:      "payment_method_investor_money",
			payMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
			err:       nil,
			expected:  36,
		},
		{
			name:      "payment_method_investor_unknown",
			payMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED,
			err:       model.ErrNotAvailablePayMethodProto,
			expected:  0,
		},
		{
			name:      "payment_method_new",
			payMethod: newPaymentMethod,
			err:       model.ErrUnexpectedPaytMethodProto,
			expected:  0,
		},
	}
	for _, test := range tests {
		req := paymentV1.PayOrderRequest{
			OrderUuid:     "any_uuid",
			UserUuid:      "any_uuid",
			PaymentMethod: test.payMethod,
		}
		res, err := api.PayOrder(a.ctx, &req)
		a.Assert().Equal(test.err, err)
		a.Assert().Len(res.GetTransactionUuid(), test.expected)
	}
}
