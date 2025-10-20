package order

import (
	"context"
	"slices"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (s *service) Create(ctx context.Context, info model.OrderInfo) (*model.OrderCreateRes, error) {
	var totalPrice float64

	parts, err := s.inventoryClient.ListParts(ctx, model.PartsFilter{
		Uuids: info.GetPartUuids(),
	})
	if err != nil {
		return nil, err
	}

	var existsPartUuids []string
	for partUuid, part := range parts {
		if part.StockQuantity <= 0 {
			continue
		}
		totalPrice += part.Price
		existsPartUuids = append(existsPartUuids, partUuid)
	}
	if len(existsPartUuids) == 0 {
		return nil, model.ErrPartsNotAvailable
	}

	slices.Sort(existsPartUuids)

	order := model.Order{
		UserUuid:   info.GetUserUuid(),
		PartUuids:  existsPartUuids,
		TotalPrice: totalPrice,
	}

	orderUuid, err := s.orderRepository.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return &model.OrderCreateRes{
		OrderUuid:  orderUuid,
		TotalPrice: order.TotalPrice,
	}, nil
}
