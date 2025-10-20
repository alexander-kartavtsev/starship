package converter

import (
	"github.com/alexander-kartavtsev/starship/order/internal/model"
	repoModel "github.com/alexander-kartavtsev/starship/order/internal/repository/model"
)

func OrderToModel(order repoModel.Order) model.Order {
	return model.Order{
		OrderUuid:       order.OrderUuid,
		UserUuid:        order.UserUuid,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUuid: order.TransactionUuid,
		PaymentMethod:   model.PaymentMethod(order.PaymentMethod),
		Status:          model.OrderStatus(order.Status),
	}
}

func OrderToRepoModel(order model.Order) repoModel.Order {
	return repoModel.Order{
		OrderUuid:       order.OrderUuid,
		UserUuid:        order.UserUuid,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUuid: order.TransactionUuid,
		PaymentMethod:   repoModel.PaymentMethod(order.PaymentMethod),
		Status:          repoModel.OrderStatus(order.Status),
	}
}
