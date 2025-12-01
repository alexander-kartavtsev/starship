package kafka

import "github.com/alexander-kartavtsev/starship/assembly/internal/model"

type AssemblyDecoder interface {
	Decode(data []byte) (model.OrderKafkaEvent, error)
}
