package order

import (
	"context"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	"github.com/alexander-kartavtsev/starship/platform/pkg/tracing"
)

func (s *service) Cancel(ctx context.Context, orderUuid string) error {
	ctx, span := tracing.StartSpan(ctx, "order.service.Cancle",
		trace.WithAttributes(
			attribute.String("order_uuid", orderUuid),
		),
	)
	defer span.End()

	order, err := s.orderRepository.Get(ctx, orderUuid)
	if err != nil {
		span.RecordError(err)
		return err
	}

	status := order.Status
	if status == model.Paid {
		span.RecordError(err)
		return model.ErrCancelPaidOrder
	}

	err = s.orderRepository.Update(
		ctx,
		orderUuid,
		model.OrderUpdateInfo{
			Status: lo.ToPtr(model.Cancelled),
		},
	)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
