package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	"github.com/alexander-kartavtsev/starship/platform/pkg/kafka"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
	eventsV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/events/v1"
)

type service struct {
	orderProducer kafka.Producer
}

func NewService(orderProducer kafka.Producer) *service {
	return &service{
		orderProducer: orderProducer,
	}
}

func (s *service) ProduceOrder(ctx context.Context, event model.OrderKafkaEvent) error {
	msg := &eventsV1.Order{
		EventUuid:       event.Uuid,
		OrderUuid:       event.OrderUuid,
		UserUuid:        event.UserUuid,
		PaymentMethod:   string(event.PaymentMethod),
		TransactionUuid: event.TransactionUuid,
		Type:            event.Type,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "Ошибка сериализации OrderRecorded", zap.Error(err))
		return err
	}

	err = s.orderProducer.Send(ctx, []byte(event.Uuid), payload)
	if err != nil {
		logger.Error(ctx, "Ошибка публикации OrderRecorded", zap.Error(err))
		return err
	}

	return nil
}
