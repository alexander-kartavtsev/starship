package order

import (
	"log"

	"github.com/go-faster/errors"
	"github.com/samber/lo"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (s *ServiceSuite) TestService_Cansel() {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "cancel_order_ok",
			err:  nil,
		},
		{
			name: "cancel_order_error",
			err:  model.ErrCancelPaidOrder,
		},
	}

	for _, test := range tests {
		log.Println(test.name)
		updateInfo := model.OrderUpdateInfo{
			Status: lo.ToPtr(model.Cancelled),
		}
		s.orderRepository.
			On("Update", s.ctx, "any_uuid", updateInfo).
			Return(test.err).
			Once()
		err := s.service.Cansel(s.ctx, "any_uuid")
		s.Assert().True(errors.Is(err, test.err))
	}
}
