package order_consumer

import (
	kafkaConverter "assembly/internal/converter/kafka"
	def "assembly/internal/service"
	"context"
	"platform/pkg/kafka"
	"platform/pkg/logger"

	"go.uber.org/zap"
)



var _ def.ConsumerService = (*service)(nil)


type service struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder kafkaConverter.OrderPaidDecoder
}

func NewService(orderPaidConsumer kafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder: orderPaidDecoder,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order orderPaidConsumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}