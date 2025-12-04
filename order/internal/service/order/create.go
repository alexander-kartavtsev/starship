package order

import (
	"context"
	"log"
	"slices"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	"github.com/alexander-kartavtsev/starship/platform/pkg/middleware/grpc"
	"github.com/alexander-kartavtsev/starship/platform/pkg/tracing"
)

func (s *service) Create(ctx context.Context, info model.OrderInfo) (*model.OrderCreateRes, error) {
	var totalPrice float64

	ctx = grpc.ForwardSessionUUIDToGRPC(ctx)
	log.Printf("ctxr: %v\n", ctx)

	// start: Span получение данных из inventory
	ctx, span := tracing.StartSpan(ctx, "inventory.list_parts",
		trace.WithAttributes(
			attribute.String("user.uuid", info.GetUserUuid()),
			attribute.StringSlice("part_uuids", info.GetPartUuids()),
		),
	)

	parts, err := s.inventoryClient.ListParts(ctx, model.PartsFilter{
		Uuids: info.GetPartUuids(),
	})
	if err != nil {
		span.RecordError(err)
		span.End()
		log.Printf("listParts err: %v\n", err)
		return nil, err
	}
	log.Printf("listParts: %v\n", parts)

	var existsPartUuids []string
	for partUuid, part := range parts {
		if part.StockQuantity <= 0 {
			continue
		}
		totalPrice += part.Price
		existsPartUuids = append(existsPartUuids, partUuid)
	}
	if len(existsPartUuids) == 0 {
		span.RecordError(err)
		span.End()
		return nil, model.ErrPartsNotAvailable
	}
	span.SetAttributes(
		attribute.StringSlice("exists_part_uuids", existsPartUuids),
	)
	span.End()
	// end: Span получение данных из inventory

	slices.Sort(existsPartUuids)

	order := model.Order{
		UserUuid:   info.GetUserUuid(),
		PartUuids:  existsPartUuids,
		TotalPrice: totalPrice,
	}

	ctx, span = tracing.StartSpan(ctx, "order.repository.Create",
		trace.WithAttributes(
			attribute.String("user.uuid", info.GetUserUuid()),
			attribute.StringSlice("part_uuids", existsPartUuids),
			attribute.Float64("total_price", totalPrice),
		),
	)
	defer span.End()

	orderUuid, err := s.orderRepository.Create(ctx, order)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &model.OrderCreateRes{
		OrderUuid:  orderUuid,
		TotalPrice: order.TotalPrice,
	}, nil
}
