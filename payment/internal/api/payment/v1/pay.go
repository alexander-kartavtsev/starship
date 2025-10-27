package v1

import (
	"context"
	"errors"

	"github.com/alexander-kartavtsev/starship/payment/internal/converter"
	"github.com/alexander-kartavtsev/starship/payment/internal/model"
	paymentV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transactionUuid, err := a.paymentService.Pay(ctx, converter.PayOrderRequestToModel(req))
	if err != nil {
		if errors.Is(err, model.ErrNotAvailablePaymentMethod) {
			return nil, model.ErrNotAvailablePayMethodProto
		} else {
			return nil, model.ErrUnexpectedPaytMethodProto
		}
	}
	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUuid,
	}, nil
}
