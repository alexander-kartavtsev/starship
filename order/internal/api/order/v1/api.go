package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	"github.com/alexander-kartavtsev/starship/order/internal/service"
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
)

type api struct {
	orderV1.UnimplementedHandler
	orderService service.OrderService
}

func NewApi(orderService service.OrderService) *api {
	return &api{
		orderService: orderService,
	}
}

func (a *api) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	status := http.StatusBadGateway
	if errors.Is(err, model.ErrCancelPaidOrder) {
		status = http.StatusConflict
	}
	if errors.Is(err, model.ErrOrderNotFound) {
		status = http.StatusNotFound
	}
	return &orderV1.GenericErrorStatusCode{
		StatusCode: status,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(status),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}
