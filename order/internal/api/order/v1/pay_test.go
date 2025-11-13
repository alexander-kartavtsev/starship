package v1

import (
	"errors"

	"github.com/google/uuid"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
)

func (a *ApiSuite) TestApi_PayOrderByUuid_Ok() {
	paymentMethod := model.Card
	orderUuid := uuid.New()
	req := orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PaymentMethod(paymentMethod),
	}

	params := orderV1.PayOrderByUuidParams{
		OrderUUID: orderUuid,
	}

	transactionUuid := "any_transaction_uuid"

	a.orderService.
		On("Pay", a.ctx, orderUuid.String(), paymentMethod).
		Return(transactionUuid, nil).
		Once()

	expRes := orderV1.PayOrderResponse{
		TransactionUUID: transactionUuid,
	}

	res, err := a.api.PayOrderByUuid(a.ctx, &req, params)

	a.Assert().Nil(err)
	a.Assert().Equal(&expRes, res)
}

func (a *ApiSuite) TestApi_PayOrderByUuid_Err() {
	paymentMethod := model.Card
	orderUuid := uuid.New()
	req := orderV1.PayOrderRequest{
		PaymentMethod: orderV1.PaymentMethod(paymentMethod),
	}

	params := orderV1.PayOrderByUuidParams{
		OrderUUID: orderUuid,
	}

	testErr := errors.New("test error")

	a.orderService.
		On("Pay", a.ctx, orderUuid.String(), paymentMethod).
		Return("", testErr).
		Once()

	res, err := a.api.PayOrderByUuid(a.ctx, &req, params)

	a.Assert().Nil(res)
	a.Assert().True(errors.Is(err, testErr))
}
