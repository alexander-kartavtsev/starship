package v1

import (
	"errors"

	"github.com/google/uuid"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
)

func (a *ApiSuite) TestApi_CreateOrder_Ok() {
	userUuid := "any_user_uuid"
	partUuid := "any_part_uuid"
	orderUuid := "any_order_uuid"
	totalPrice := 123.45
	partUuids := []string{partUuid}
	modelOrderInfo := model.OrderInfo{
		UserUuid:  userUuid,
		PartUuids: partUuids,
	}
	orderCreateRes := model.OrderCreateRes{
		OrderUuid:  orderUuid,
		TotalPrice: totalPrice,
	}
	a.orderService.
		On("Create", a.ctx, modelOrderInfo).
		Return(&orderCreateRes, nil).
		Once()
	request := orderV1.CreateOrderRequest{
		UserUUID:  userUuid,
		PartUuids: partUuids,
	}
	response := orderV1.CreateOrderResponse{
		OrderUUID:  orderUuid,
		TotalPrice: totalPrice,
	}
	params := orderV1.CreateOrderParams{
		XSessionUUID: uuid.New(),
	}
	res, err := a.api.CreateOrder(a.ctx, &request, params)
	a.Assert().Equal(&response, res)
	a.Assert().Nil(err)
}

func (a *ApiSuite) TestApi_CreateOrder_Err() {
	userUuid := "any_user_uuid"
	partUuid := "any_part_uuid"
	partUuids := []string{partUuid}
	modelOrderInfo := model.OrderInfo{
		UserUuid:  userUuid,
		PartUuids: partUuids,
	}
	request := orderV1.CreateOrderRequest{
		UserUUID:  userUuid,
		PartUuids: partUuids,
	}
	params := orderV1.CreateOrderParams{
		XSessionUUID: uuid.New(),
	}

	a.orderService.
		On("Create", a.ctx, modelOrderInfo).
		Return(nil, model.ErrOrderNotFound).
		Once()
	res, err := a.api.CreateOrder(a.ctx, &request, params)
	a.Assert().Nil(res)
	a.Assert().True(errors.Is(err, model.ErrOrderNotFound))
}
