package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/alexander-kartavtsev/starship/payment/internal/converter"
	"github.com/alexander-kartavtsev/starship/payment/internal/model"
	paymentV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transactionUuid, err := a.paymentService.Pay(ctx, converter.PayOrderRequestToModel(req))
	if err != nil {
		if errors.Is(err, model.ErrNotAvailablePaymentMethod) {
			return nil, status.Errorf(codes.InvalidArgument, "Способ оплаты не поддерживается")
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "Неправильный метод оплаты")
		}
	}
	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUuid,
	}, nil
}
