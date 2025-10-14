package v1

import (
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
