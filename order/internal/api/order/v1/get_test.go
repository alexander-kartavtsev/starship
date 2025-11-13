package v1

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
)

func (a *ApiSuite) TestApi_GetOrderByUuid_Ok() {
	orderUuid := uuid.New()
	userUuid := "any_user_uuid"
	partUuids := []string{"any_part_uuid_1", "any_part_uuid_2"}
	transactionUuid := "any_transaction_uuid"
	modelOrder := model.Order{
		OrderUuid:       orderUuid.String(),
		UserUuid:        userUuid,
		PartUuids:       partUuids,
		TotalPrice:      123.45,
		TransactionUuid: transactionUuid,
		PaymentMethod:   model.Card,
		Status:          model.PendingPayment,
	}
	protoOrder := orderV1.OrderDto{
		OrderUUID:       orderUuid.String(),
		UserUUID:        userUuid,
		PartUuids:       partUuids,
		TotalPrice:      123.45,
		TransactionUUID: transactionUuid,
		PaymentMethod:   orderV1.PaymentMethod(model.Card),
		Status:          orderV1.OrderStatus(model.PendingPayment),
	}
	params := orderV1.GetOrderByUuidParams{
		OrderUUID: orderUuid,
	}
	a.orderService.
		On("Get", a.ctx, orderUuid.String()).
		Return(modelOrder, nil).
		Once()
	res, err := a.api.GetOrderByUuid(a.ctx, params)
	a.Assert().Equal(&protoOrder, res)
	a.Assert().Nil(err)
}

func (a *ApiSuite) TestApi_GetOrderByUuid_Err() {
	orderUuid := uuid.New()
	modelOrder := model.Order{}
	params := orderV1.GetOrderByUuidParams{
		OrderUUID: orderUuid,
	}
	serviceErr := model.ErrOrderNotFound
	a.orderService.
		On("Get", a.ctx, orderUuid.String()).
		Return(modelOrder, serviceErr).
		Once()
	response := orderV1.NotFoundError{
		Code:    http.StatusNotFound,
		Message: "Заказ с id = '" + params.OrderUUID.String() + "' не найден",
	}
	res, err := a.api.GetOrderByUuid(a.ctx, params)
	a.Assert().Equal(&response, res)
	a.Assert().True(errors.Is(err, model.ErrOrderNotFound))
}
