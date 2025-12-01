package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/alexander-kartavtsev/starship/notification/internal/converter/kafka"
	def "github.com/alexander-kartavtsev/starship/notification/internal/service"
	"github.com/alexander-kartavtsev/starship/platform/pkg/kafka"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
)

var _ def.OrderAssembledConsumerService = (*service)(nil)

type service struct {
	consumer        kafka.Consumer
	decoder         kafkaConverter.OrderAssembledDecoder
	telegramService def.TelegramService
}

func NewService(consumer kafka.Consumer, decoder kafkaConverter.OrderAssembledDecoder, telegramService def.TelegramService) *service {
	return &service{
		consumer:        consumer,
		decoder:         decoder,
		telegramService: telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderConsumer service")

	err := s.consumer.Consume(ctx, s.OrderAssembledHandler)
	if err != nil {
		logger.Error(ctx, "Consume from OrderAssembled topic error", zap.Error(err))
		return err
	}

	return nil
}
