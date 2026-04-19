package service

import (
	"context"

	"assembly/internal/model"
)

type AssemblyConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type AssemblyProducerService interface {
	ProduceOrderAssembled(ctx context.Context, event model.OrderAssembledEvent) error
}
