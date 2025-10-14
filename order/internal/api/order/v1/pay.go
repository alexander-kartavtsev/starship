package v1

import (
	"context"
	"github.com/alexander-kartavtsev/starship/order/internal/model"
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrderByUuid(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderByUuidParams) (orderV1.PayOrderByUuidRes, error) {
	transactionUuid, err := a.orderService.Pay(
		ctx,
		params.OrderUUID.String(),
		model.PaymentMethod(req.GetPaymentMethod()),
	)
	if err != nil {
		return nil, err
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUuid,
	}, nil
}
