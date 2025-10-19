package order

import (
	"errors"
	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (s *ServiceSuite) TestService_CreateOk() {
	partUuids := []string{"any_part_uuid"}
	userUuid := "any_user_uuid"
	price := 123.45
	modelOrderInfo := model.OrderInfo{
		UserUuid:  userUuid,
		PartUuids: partUuids,
	}
	testParts := map[string]model.Part{
		"any_part_uuid": model.Part{
			Uuid:          "any_part_uuid",
			Price:         price,
			StockQuantity: 12,
		},
	}
	testOrder := model.Order{
		UserUuid:   userUuid,
		PartUuids:  partUuids,
		TotalPrice: price,
	}
	createRes := &model.OrderCreateRes{
		OrderUuid:  "any_order_uuid",
		TotalPrice: price,
	}

	s.inventoryClient.
		On("ListParts", s.ctx, model.PartsFilter{Uuids: partUuids}).
		Return(testParts, nil).
		Once()

	s.orderRepository.
		On("Create", s.ctx, testOrder).
		Return("any_order_uuid", nil).
		Once()

	res, err := s.service.Create(s.ctx, modelOrderInfo)
	s.Assert().Nil(err)
	s.Assert().Equal(res, createRes)
}

func (s *ServiceSuite) TestService_CreateSemiOk() {
	partUuids := []string{"any_part_uuid_01", "any_part_uuid_02", "any_part_uuid_03"}
	userUuid := "any_user_uuid"
	price1 := 125.0
	price2 := 250.0
	modelOrderInfo := model.OrderInfo{
		UserUuid:  userUuid,
		PartUuids: partUuids,
	}
	testParts := map[string]model.Part{
		"any_part_uuid_01": model.Part{
			Uuid:          "any_part_uuid_01",
			Price:         price1,
			StockQuantity: 12,
		},
		"any_part_uuid_02": model.Part{
			Uuid:          "any_part_uuid_02",
			Price:         price2,
			StockQuantity: 12,
		},
	}
	testOrder := model.Order{
		UserUuid:   userUuid,
		PartUuids:  []string{"any_part_uuid_01", "any_part_uuid_02"},
		TotalPrice: price1 + price2,
	}
	createRes := &model.OrderCreateRes{
		OrderUuid:  "any_order_uuid",
		TotalPrice: price1 + price2,
	}

	s.inventoryClient.
		On("ListParts", s.ctx, model.PartsFilter{Uuids: partUuids}).
		Return(testParts, nil).
		Once()

	s.orderRepository.
		On("Create", s.ctx, testOrder).
		Return("any_order_uuid", nil).
		Once()

	res, err := s.service.Create(s.ctx, modelOrderInfo)
	s.Assert().Nil(err)
	s.Assert().Equal(res, createRes)
}

func (s *ServiceSuite) TestService_CreatePartsNotAvailability() {
	partUuids := []string{"any_part_uuid_01", "any_part_uuid_02"}
	userUuid := "any_user_uuid"
	modelOrderInfo := model.OrderInfo{
		UserUuid:  userUuid,
		PartUuids: partUuids,
	}

	tests := []struct {
		parts map[string]model.Part
	}{
		{
			parts: map[string]model.Part{},
		},
		{
			parts: map[string]model.Part{
				"any_part_uuid_01": model.Part{
					Uuid:          "any_part_uuid_01",
					Price:         123.45,
					StockQuantity: 0,
				},
				"any_part_uuid_02": model.Part{
					Uuid:          "any_part_uuid_02",
					Price:         123.45,
					StockQuantity: 0,
				},
			},
		},
	}

	for _, test := range tests {
		s.inventoryClient.
			On("ListParts", s.ctx, model.PartsFilter{Uuids: partUuids}).
			Return(test.parts, nil).
			Once()

		res, err := s.service.Create(s.ctx, modelOrderInfo)
		s.Assert().Nil(res)
		s.Assert().True(errors.Is(err, model.ErrPartsNotAvailable))
	}
}

func (s *ServiceSuite) TestService_CreateErrClient() {
	partUuids := []string{"any_part_uuid"}
	modelOrderInfo := model.OrderInfo{
		UserUuid:  "any_user_uuid",
		PartUuids: partUuids,
	}

	testError := errors.New("any test error")

	s.inventoryClient.
		On("ListParts", s.ctx, model.PartsFilter{Uuids: partUuids}).
		Return(nil, testError).
		Once()

	res, err := s.service.Create(s.ctx, modelOrderInfo)
	s.Assert().True(errors.Is(err, testError))
	s.Assert().Nil(res)
}

func (s *ServiceSuite) TestService_CreateErrRepository() {
	partUuids := []string{"any_part_uuid"}
	userUuid := "any_user_uuid"
	price := 123.45
	modelOrderInfo := model.OrderInfo{
		UserUuid:  userUuid,
		PartUuids: partUuids,
	}
	testParts := map[string]model.Part{
		"any_part_uuid": model.Part{
			Uuid:          "any_part_uuid",
			Price:         price,
			StockQuantity: 12,
		},
	}
	testOrder := model.Order{
		UserUuid:   userUuid,
		PartUuids:  partUuids,
		TotalPrice: price,
	}
	testError := errors.New("any test error")

	s.inventoryClient.
		On("ListParts", s.ctx, model.PartsFilter{Uuids: partUuids}).
		Return(testParts, nil).
		Once()
	s.orderRepository.
		On("Create", s.ctx, testOrder).
		Return("", testError).
		Once()

	res, err := s.service.Create(s.ctx, modelOrderInfo)
	s.Assert().True(errors.Is(err, testError))
	s.Assert().Nil(res)
}
