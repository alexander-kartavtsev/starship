package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/alexander-kartavtsev/starship/order/internal/converter/kafka"
	"github.com/alexander-kartavtsev/starship/platform/pkg/kafka"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
)

type service struct {
	orderConsumer kafka.Consumer
	orderDecoder  kafkaConverter.OrderDecoder
}

func NewService(orderConsumer kafka.Consumer, orderDecoder kafkaConverter.OrderDecoder) *service {
	return &service{
		orderConsumer: orderConsumer,
		orderDecoder:  orderDecoder,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderConsumer service")

	err := s.orderConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order topic error", zap.Error(err))
		return err
	}

	return nil
}
