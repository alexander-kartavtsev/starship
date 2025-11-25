package assembly_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/alexander-kartavtsev/starship/assembly/internal/converter/kafka"
	assemblyService "github.com/alexander-kartavtsev/starship/assembly/internal/service"
	"github.com/alexander-kartavtsev/starship/platform/pkg/kafka"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
)

type service struct {
	assemblyConsumer        kafka.Consumer
	assemblyDecoder         kafkaConverter.AssemblyDecoder
	assemblyProducerService assemblyService.ProducerService
}

func NewService(consumer kafka.Consumer, decoder kafkaConverter.AssemblyDecoder, producerService assemblyService.ProducerService) *service {
	return &service{
		assemblyConsumer:        consumer,
		assemblyDecoder:         decoder,
		assemblyProducerService: producerService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderConsumer service")

	err := s.assemblyConsumer.Consume(ctx, s.AssemblyHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order topic error", zap.Error(err))
		return err
	}

	return nil
}
