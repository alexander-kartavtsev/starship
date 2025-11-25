package decoder

import (
	"google.golang.org/protobuf/proto"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	eventsV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewAssemblyDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.ShipAssembledKafkaEvent, error) {
	var pb eventsV1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledKafkaEvent{}, err
	}

	return model.ShipAssembledKafkaEvent{
		EventUuid:    pb.EventUuid,
		OrderUuid:    pb.OrderUuid,
		UserUuid:     pb.UserUuid,
		BuildTimeSec: pb.BuildTimeSec,
	}, nil
}
