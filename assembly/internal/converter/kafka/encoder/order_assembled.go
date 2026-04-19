package encoder

import (
	"fmt"

	"assembly/internal/model"
	"google.golang.org/protobuf/proto"
	eventsV1 "shared/pkg/proto/events/v1"
)

type encoder struct{}

func NewOrderAssembledEncoder() *encoder {
	return &encoder{}
}

func (e *encoder) Encode(event model.OrderAssembledEvent) ([]byte, error) {
	pb := &eventsV1.OrderAssembled{
		EventUuid:    event.EventUuid.String(),
		OrderUuid:    event.OrderUuid.String(),
		UserUuid:     event.UserUuid.String(),
		BuildTimeSec: event.BuildTimeSec,
	}

	data, err := proto.Marshal(pb)
	if err != nil {
		return nil, fmt.Errorf("proto marshal order assembled: %w", err)
	}
	return data, nil
}
