package v1

import (
	"context"
	"net/http"

	"github.com/alexander-kartavtsev/starship/order/internal/converter"
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUuid(ctx context.Context, params orderV1.GetOrderByUuidParams) (orderV1.GetOrderByUuidRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID.String())
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "Заказ с id = '" + params.OrderUUID.String() + "' не найден",
		}, err
	}
	return converter.OrderToApi(order), nil
}
