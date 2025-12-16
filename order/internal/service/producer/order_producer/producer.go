package orderproducer

import (
	"context"
	"order/internal/model"
	"platform/pkg/kafka"
	"platform/pkg/logger"
	eventsV1 "shared/pkg/proto/events/v1"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	def "order/internal/service"
)

var _ def.OrderProducerService = (*service)(nil)


type service struct {
	orderPaidProducer kafka.Producer
}

func NewService(orderPaidProducer kafka.Producer) *service {
	return &service{
		orderPaidProducer: orderPaidProducer,
	}
}

func (s *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	msg := &eventsV1.OrderPaid{
		EventUuid:       event.EventUuid.String(),
		OrderUuid:       event.OrderUuid.String(),
		UserUuid:        event.UserUuid.String(),
		PaymentMethod:   eventsV1.PaymentMethod(event.PaymentMethod),
		TransactionUuid: event.TransactionUuid.String(),
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderPaid", zap.Error(err))
		return err
	}

	err = s.orderPaidProducer.Send(ctx, []byte(event.EventUuid.String()), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish OrderPaid", zap.Error(err))
		return err
	}

	return nil

}
