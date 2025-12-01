package decoder

import (
	"google.golang.org/protobuf/proto"

	"github.com/alexander-kartavtsev/starship/notification/internal/model"
	eventsV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/events/v1"
)

type orderPaidDecoder struct{}

func NewOrderPaidDecoder() *orderPaidDecoder {
	return &orderPaidDecoder{}
}

func (d *orderPaidDecoder) Decode(data []byte) (model.OrderKafkaEvent, error) {
	var pb eventsV1.Order
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderKafkaEvent{}, err
	}

	return model.OrderKafkaEvent{
		Uuid:            pb.EventUuid,
		OrderUuid:       pb.OrderUuid,
		UserUuid:        pb.UserUuid,
		PaymentMethod:   model.PaymentMethod(pb.PaymentMethod),
		TransactionUuid: pb.TransactionUuid,
		Type:            pb.Type,
	}, nil
}
