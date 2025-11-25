package service

import (
	"context"

	"github.com/alexander-kartavtsev/starship/notification/internal/model"
)

type OrderAssembledConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type TelegramService interface {
	SendAssembledNotification(ctx context.Context, event model.ShipAssembledKafkaEvent) error
	SendPaidNotification(ctx context.Context, event model.OrderKafkaEvent) error
}
