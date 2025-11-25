package decoder

import (
	"google.golang.org/protobuf/proto"

	"github.com/alexander-kartavtsev/starship/notification/internal/model"
	eventsV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/events/v1"
)

type orderAssembledDecoder struct{}

func NewOrderAssembledDecoder() *orderAssembledDecoder {
	return &orderAssembledDecoder{}
}

func (d *orderAssembledDecoder) Decode(data []byte) (model.ShipAssembledKafkaEvent, error) {
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
