package order

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	"github.com/alexander-kartavtsev/starship/platform/pkg/tracing"
)

func (s *service) Get(ctx context.Context, uuid string) (model.Order, error) {
	ctx, span := tracing.StartSpan(ctx, "order.service.Get",
		trace.WithAttributes(
			attribute.String("order_uuid", uuid),
		),
	)
	defer span.End()

	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		span.RecordError(err)
		return model.Order{}, err
	}
	return order, nil
}
