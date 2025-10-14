package v1

import (
	"context"
	"net/http"

	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrderById(ctx context.Context, params orderV1.CancelOrderByIdParams) (orderV1.CancelOrderByIdRes, error) {
	err := a.orderService.Cansel(ctx, params.OrderUUID.String())
	if err != nil {
		return nil, a.NewError(ctx, err)
	}

	return &orderV1.CancelOrderByIdNoContent{
		Code:    http.StatusNoContent,
		Message: "Заказ отменен",
	}, nil
}
