package order

import (
	"log"

	"github.com/go-faster/errors"
	"github.com/samber/lo"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	"github.com/alexander-kartavtsev/starship/platform/pkg/tracing"
)

func (s *ServiceSuite) TestService_Cancel() {
	tests := []struct {
		name   string
		order  model.Order
		getErr error
		err    error
	}{
		{
			name:   "cancel_order_ok",
			order:  model.Order{Status: model.PendingPayment},
			getErr: nil,
			err:    nil,
		},
		{
			name:   "cancel_order_cansel_error",
			order:  model.Order{Status: model.PendingPayment},
			getErr: nil,
			err:    model.ErrCancelPaidOrder,
		},
	}

	ctx, span := tracing.StartSpan(s.ctx, "order.service.Cancle")
	span.End()

	for _, test := range tests {
		log.Println(test.name)
		updateInfo := model.OrderUpdateInfo{
			Status: lo.ToPtr(model.Cancelled),
		}
		s.orderRepository.
			On("Get", ctx, "any_uuid").
			Return(test.order, test.getErr).
			Once()
		s.orderRepository.
			On("Update", ctx, "any_uuid", updateInfo).
			Return(test.err).
			Once()
		err := s.service.Cancel(s.ctx, "any_uuid")
		s.Assert().True(errors.Is(err, test.err))
	}
}

func (s *ServiceSuite) TestService_Cansel_Not() {
	tests := []struct {
		name   string
		order  model.Order
		getErr error
		err    error
	}{
		{
			name:   "cancel_order_not_found_error",
			order:  model.Order{},
			getErr: model.ErrOrderNotFound,
			err:    model.ErrOrderNotFound,
		},
		{
			name:   "cancel_order_not_found_error",
			order:  model.Order{Status: model.Paid},
			getErr: nil,
			err:    model.ErrCancelPaidOrder,
		},
	}

	ctx, span := tracing.StartSpan(s.ctx, "")
	span.End()

	for _, test := range tests {
		log.Println(test.name)
		s.orderRepository.
			On("Get", ctx, "any_uuid").
			Return(test.order, test.getErr).
			Once()
		err := s.service.Cancel(s.ctx, "any_uuid")
		s.Assert().True(errors.Is(err, test.err))
	}
}
