package service

import "context"

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}


type ProducerService interface {
	RunProducer(ctx context.Context) error
}
