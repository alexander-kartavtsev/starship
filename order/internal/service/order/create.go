package order

import (
	"context"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (s *service) Create(ctx context.Context, info model.OrderInfo) (*model.OrderCreateRes, error) {
	var totalPrice float64
	// log.Printf("info (параметр в order.service.Create): %v\n", info)

	parts, err := s.inventoryClient.ListParts(ctx, model.PartsFilter{
		Uuids: info.GetPartUuids(),
	})
	if err != nil {
		// log.Printf("parts (получили из gRPC в order.service.Create): %v\n", parts)
		// log.Printf("err (получили из gRPC в order.service.Create): %v\n", err)
		return nil, err
	}
	// log.Printf("parts (получили из gRPC в order.service.Create): %v\n", parts)
	var existsPartUuids []string
	for partUuid, part := range parts {
		totalPrice += part.Price
		existsPartUuids = append(existsPartUuids, partUuid)
	}

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
