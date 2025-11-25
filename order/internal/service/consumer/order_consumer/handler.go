package order_consumer

import (
	"context"
	"log"

	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	"github.com/alexander-kartavtsev/starship/platform/pkg/kafka"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode Order", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("Uuid", event.EventUuid),
		zap.String("OrderUuid", event.OrderUuid),
		zap.String("UserUuid", event.UserUuid),
		zap.Int64("BuildTimeSec", event.BuildTimeSec),
	)

	order, err := s.orderRepository.Get(ctx, event.OrderUuid)
	if err != nil {
		return err
	}

	log.Printf("Заказ %s имеет статус %s", order.OrderUuid, order.Status)

	err = s.orderRepository.Update(
		ctx,
		order.OrderUuid,
		model.OrderUpdateInfo{
			Status: lo.ToPtr(model.Assembled),
		},
	)
	if err != nil {
		return err
	}

	order, err = s.orderRepository.Get(ctx, order.OrderUuid)
	if err != nil {
		return err
	}

	log.Printf("Заказу %s установлен статус %s", order.OrderUuid, order.Status)

	return nil
}
