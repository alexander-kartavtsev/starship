package service

import (
	"context"

	"github.com/alexander-kartavtsev/starship/assembly/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ProducerService interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembledKafkaEvent) error
}
