package kafka

import "assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}
