package kafka

import "assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}

type OrderAssembledEncoder interface {
	Encode(event model.OrderAssembledEvent) ([]byte, error)
}

type DLQEventEncoder interface {
	Encode(event model.DeadLetterEvent) ([]byte, error)
}
