package kafka

import "github.com/alexander-kartavtsev/starship/order/internal/model"

type OrderDecoder interface {
	Decode(data []byte) (model.ShipAssembledKafkaEvent, error)
}
