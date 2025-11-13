package v1

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
)

func (a *ApiSuite) TestApi_CancelOrderById_Ok() {
	testUuid := uuid.New()
	expectedRes := &orderV1.CancelOrderByIdNoContent{
		Code:    http.StatusNoContent,
		Message: "Заказ отменен",
	}

	a.orderService.
		On("Cancel", a.ctx, testUuid.String()).
		Return(nil).
		Once()
	res, err := a.api.CancelOrderById(a.ctx, orderV1.CancelOrderByIdParams{
		OrderUUID: testUuid,
	})
	a.Assert().True(errors.Is(err, nil))
	a.Assert().Equal(expectedRes, res)
}

func (a *ApiSuite) TestApi_CancelOrderById_Err() {
	testUuid := uuid.New()

	a.orderService.
		On("Cancel", a.ctx, testUuid.String()).
		Return(model.ErrOrderNotFound).
		Once()

	res, err := a.api.CancelOrderById(a.ctx, orderV1.CancelOrderByIdParams{
		OrderUUID: testUuid,
	})
	a.Assert().NotNil(err)
	a.Assert().Nil(res)
}

func (a *ApiSuite) TestApi_CancelOrderById_Err_Conflict() {
	testUuid := uuid.New()

	a.orderService.
		On("Cancel", a.ctx, testUuid.String()).
		Return(model.ErrCancelPaidOrder).
		Once()

	res, err := a.api.CancelOrderById(a.ctx, orderV1.CancelOrderByIdParams{
		OrderUUID: testUuid,
	})
	a.Assert().NotNil(err)
	a.Assert().Nil(res)
}
