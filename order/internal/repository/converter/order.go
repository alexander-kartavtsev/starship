package converter

import (
	"github.com/alexander-kartavtsev/starship/order/internal/model"
	repoModel "github.com/alexander-kartavtsev/starship/order/internal/repository/model"
)

func OrderInfoToRepoModel(info model.OrderInfo) repoModel.OrderInfo {
	return repoModel.OrderInfo{}
}

func OrderToModel(order repoModel.Order) model.Order {
	return model.Order{}
}

func OrderInfoToModel(userUuid string, partUuids []string) model.OrderInfo {
	return model.OrderInfo{}
}
