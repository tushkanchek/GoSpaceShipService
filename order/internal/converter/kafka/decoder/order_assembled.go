package decoder

import (
	"fmt"
	def "order/internal/converter/kafka"
	"order/internal/model"
	eventsV1 "shared/pkg/proto/events/v1"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

var _ def.OrderAssembledDecoder = (*decoder)(nil)

type decoder struct {
}

func NewDecoder() *decoder {
	return &decoder{}
}


func (d *decoder) Decode(data []byte) (model.OrderAssembledEvent, error) {
	var pb eventsV1.OrderAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderAssembledEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}
	EventUuid, err := parseUUID(pb.GetEventUuid(), "EventUuid")
	if err != nil {
		return model.OrderAssembledEvent{}, err
	}
	OrderUuid, err := parseUUID(pb.GetOrderUuid(), "OrderUuid")
	if err != nil {
		return model.OrderAssembledEvent{}, err
	}
	UserUuid, err := parseUUID(pb.GetUserUuid(), "UserUuid")
	if err != nil {
		return model.OrderAssembledEvent{}, err
	}	

	return model.OrderAssembledEvent{
		EventUuid:  EventUuid,
		OrderUuid:  OrderUuid,
		UserUuid:   UserUuid,
		BuildTimeSec: pb.GetBuildTimeSec(),
	}, nil
}

func parseUUID(s string, fieldName string) (uuid.UUID, error) {
    id, err := uuid.Parse(s)
    if err != nil {
        return uuid.Nil, fmt.Errorf("invalid %s: %w", fieldName, err)
    }
    return id, nil
}