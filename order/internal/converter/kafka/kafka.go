package kafka

import (
	"order/internal/model"
)




type OrderAssembledDecoder interface {
	Decode(data []byte) (model.OrderAssembledEvent, error)
}