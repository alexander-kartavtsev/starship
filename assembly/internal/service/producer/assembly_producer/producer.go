package assembly_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/alexander-kartavtsev/starship/assembly/internal/model"
	"github.com/alexander-kartavtsev/starship/platform/pkg/kafka"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
	eventsV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/events/v1"
)

type service struct {
	assemblyProducer kafka.Producer
}

func NewService(assemblyProducer kafka.Producer) *service {
	return &service{
		assemblyProducer: assemblyProducer,
	}
}

func (s *service) ProduceShipAssembled(ctx context.Context, event model.ShipAssembledKafkaEvent) error {
	msg := &eventsV1.ShipAssembled{
		EventUuid:    event.EventUuid,
		OrderUuid:    event.OrderUuid,
		UserUuid:     event.UserUuid,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "Ошибка сериализации Order", zap.Error(err))
		return err
	}

	err = s.assemblyProducer.Send(ctx, []byte(event.EventUuid), payload)
	if err != nil {
		logger.Error(ctx, "Ошибка публикации Order в ShipAssembled", zap.Error(err))
		return err
	}

	return nil
}
