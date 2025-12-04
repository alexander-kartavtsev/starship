package v1

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	"github.com/alexander-kartavtsev/starship/platform/pkg/tracing"
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest, params orderV1.CreateOrderParams) (orderV1.CreateOrderRes, error) {
	ctx, span := tracing.StartSpan(ctx, "order.v1.CreateOrder",
		trace.WithAttributes(
			attribute.String("user.uuid", req.GetUserUUID()),
			attribute.StringSlice("part_uuids", req.GetPartUuids()),
		),
	)
	defer span.End()

	res, err := a.orderService.Create(
		ctx,
		model.OrderInfo{
			UserUuid:  req.GetUserUUID(),
			PartUuids: req.GetPartUuids(),
		},
	)
	if err != nil {
		span.RecordError(err)
		span.End()
		return nil, err
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  res.OrderUuid,
		TotalPrice: res.TotalPrice,
	}, nil
}
