package order

import (
	"errors"
	"log"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	"github.com/alexander-kartavtsev/starship/platform/pkg/tracing"
)

func (s *ServiceSuite) TestService_Get() {
	tests := []struct {
		name      string
		orderRepo model.Order
		orderRes  model.Order
		err       error
	}{
		{
			name:      "order_service_get_ok",
			orderRepo: model.Order{OrderUuid: "anu_order_uuid"},
			orderRes:  model.Order{OrderUuid: "anu_order_uuid"},
			err:       nil,
		},
		{
			name:      "order_service_get_ok",
			orderRepo: model.Order{OrderUuid: "anu_order_uuid"},
			orderRes:  model.Order{},
			err:       errors.New("test_error"),
		},
	}

	ctx, span := tracing.StartSpan(s.ctx, "order.service.Get")
	span.End()

	for _, test := range tests {
		log.Println(test.name)
		s.orderRepository.
			On("Get", ctx, "anu_order_uuid").
			Return(test.orderRepo, test.err).
			Once()
		res, err := s.service.Get(s.ctx, "anu_order_uuid")
		s.Assert().True(errors.Is(err, test.err))
		s.Assert().Equal(res, test.orderRes)
	}
}
