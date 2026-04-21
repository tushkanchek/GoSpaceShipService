package orderconsumer

import (
	"context"
	kafkaDecoder "order/internal/converter/kafka"
	def "order/internal/service"
	"platform/pkg/kafka"
	"platform/pkg/logger"

	"go.uber.org/zap"
)

var _ def.OrderConsumerService = (*service)(nil)


type service struct {
	orderAssembledConsumer	  kafka.Consumer
	orderAssembledDecoder	  kafkaDecoder.OrderAssembledDecoder
	orderService			  def.OrderService
}

func NewService(orderAssembledConsumer kafka.Consumer, orderAssembledDecoder kafkaDecoder.OrderAssembledDecoder, orderService def.OrderService) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		orderService: orderService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order orderAssembledConsumer service")

	done := make(chan error, 1)
	go func() {
		done <- s.orderAssembledConsumer.Consume(ctx, s.OrderHandler)
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Consumer shutdown signal received")
		return ctx.Err()
	case err := <-done:
		if err != nil {
			logger.Error(ctx, "Consume from order.assembled topic error", zap.Error(err))
		}
		return err
	}
}