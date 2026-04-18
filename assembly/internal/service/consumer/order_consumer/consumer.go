package order_consumer

import (
	"context"

	kafkaConverter "assembly/internal/converter/kafka"
	def "assembly/internal/service"
	"go.uber.org/zap"
	"platform/pkg/kafka"
	"platform/pkg/logger"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder
}

func NewService(orderPaidConsumer kafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder:  orderPaidDecoder,
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
