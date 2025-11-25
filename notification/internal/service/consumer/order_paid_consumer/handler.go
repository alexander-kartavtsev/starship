package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/alexander-kartavtsev/starship/platform/pkg/kafka"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.decoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid in notification service", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("Uuid", event.Uuid),
		zap.String("OrderUuid", event.OrderUuid),
		zap.String("UserUuid", event.UserUuid),
		zap.String("PaymentMethod", string(event.PaymentMethod)),
		zap.String("TransactionUuid", event.TransactionUuid),
		zap.String("Type", event.Type),
	)

	// Отправляем уведомление в Telegram
	if err := s.telegramService.SendPaidNotification(ctx, event); err != nil {
		return err
	}

	return nil
}
