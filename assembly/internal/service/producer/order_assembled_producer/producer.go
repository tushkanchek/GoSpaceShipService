package order_assembled_producer

import (
	"context"
	//"time"

	kafkaConverter "assembly/internal/converter/kafka"
	"assembly/internal/model"
	def "assembly/internal/service"
	"go.uber.org/zap"
	"platform/pkg/kafka"
	"platform/pkg/logger"
)

var _ def.AssemblyProducerService = (*service)(nil)

type service struct {
	orderAssembledProducer kafka.Producer
	orderAssembledEncoder  kafkaConverter.OrderAssembledEncoder
}

func NewService(orderAssembledProducer kafka.Producer, orderAssembledEncoder kafkaConverter.OrderAssembledEncoder) *service {
	return &service{
		orderAssembledProducer: orderAssembledProducer,
		orderAssembledEncoder:  orderAssembledEncoder,
	}
}

func (s *service) ProduceOrderAssembled(ctx context.Context, event model.OrderAssembledEvent) error {
	assembly, err := s.orderAssembledEncoder.Encode(event)
	if err != nil {
		logger.Error(ctx, "Failed to encode OrderAssembledEvent", zap.Error(err))
		return err
	}

	// timer := time.NewTimer(10 * time.Second)
	// defer timer.Stop()

	// select {
	// case <-timer.C:
		err = s.orderAssembledProducer.Send(ctx, []byte(event.EventUuid.String()), assembly)
		if err != nil {
			logger.Error(ctx, "Failed to publish OrderAssembled", zap.Error(err))
			return err
		}
		return nil
	// case <-ctx.Done():
	// 	logger.Info(ctx, "stopping assembly wait due to context cancellation")
	// 	return ctx.Err()
	// }
}
