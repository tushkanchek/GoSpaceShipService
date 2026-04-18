package decoder

import (
	"fmt"

	"assembly/internal/model"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	eventsV1 "shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderPaidDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsV1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.OrderPaidEvent{
		EventUuid:       uuid.MustParse(pb.GetEventUuid()),
		OrderUuid:       uuid.MustParse(pb.OrderUuid),
		UserUuid:        uuid.MustParse(pb.UserUuid),
		PaymentMethod:   pb.PaymentMethod.String(),
		TransactionUuid: uuid.MustParse(pb.TransactionUuid),
	}, nil
}
