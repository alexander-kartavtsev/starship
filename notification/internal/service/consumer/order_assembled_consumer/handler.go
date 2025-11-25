package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/alexander-kartavtsev/starship/platform/pkg/kafka"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
)

func (s *service) OrderAssembledHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.decoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderAssembled in notification service", zap.Error(err))
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

	// Отправляем уведомление в Telegram
	if err := s.telegramService.SendAssembledNotification(ctx, event); err != nil {
		return err
	}

	return nil
}
