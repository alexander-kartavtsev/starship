package kafka

import "github.com/alexander-kartavtsev/starship/notification/internal/model"

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembledKafkaEvent, error)
}

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderKafkaEvent, error)
}
