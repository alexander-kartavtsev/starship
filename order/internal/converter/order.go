package converter

import (
	"github.com/alexander-kartavtsev/starship/order/internal/model"
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
)

func OrderToApi(order model.Order) *orderV1.OrderDto {
	return &orderV1.OrderDto{
		OrderUUID:       order.OrderUuid,
		UserUUID:        order.UserUuid,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUuid,
		PaymentMethod:   orderV1.PaymentMethod(order.PaymentMethod),
		Status:          orderV1.OrderStatus(order.Status),
	}
}
