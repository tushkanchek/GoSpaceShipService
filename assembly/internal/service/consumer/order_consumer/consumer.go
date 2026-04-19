package order_consumer

import (
	"context"

	kafkaConverter "assembly/internal/converter/kafka"
	def "assembly/internal/service"
	"platform/pkg/kafka"
	"platform/pkg/logger"

	"go.uber.org/zap"
)

var _ def.AssemblyConsumerService = (*service)(nil)

type service struct {
	orderPaidConsumer      kafka.Consumer
	orderPaidDecoder       kafkaConverter.OrderPaidDecoder
	orderAssembledProducer def.AssemblyProducerService
	dlqProducer            kafka.Producer
	dlqEncoder             kafkaConverter.DLQEventEncoder
}

func NewService(
	orderPaidConsumer kafka.Consumer,
	orderPaidDecoder kafkaConverter.OrderPaidDecoder,
	orderAssembledProducer def.AssemblyProducerService,
	dlqProducer kafka.Producer,
	dlqEncoder kafkaConverter.DLQEventEncoder,
) *service {
	return &service{
		orderPaidConsumer:      orderPaidConsumer,
		orderPaidDecoder:       orderPaidDecoder,
		orderAssembledProducer: orderAssembledProducer,
		dlqProducer:            dlqProducer,
		dlqEncoder:             dlqEncoder,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order orderPaidConsumer service")

	done := make(chan error, 1)
	go func() {
		done <- s.orderPaidConsumer.Consume(ctx, s.OrderHandler)
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Consumer shutdown signal received")
		return ctx.Err()
	case err := <-done:
		if err != nil {
			logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		}
		return err
	}
}
