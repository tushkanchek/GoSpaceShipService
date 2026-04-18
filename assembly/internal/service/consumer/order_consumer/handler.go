package order_consumer

import (
	"context"

	"go.uber.org/zap"
	"platform/pkg/kafka"
	"platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUuid.String()),
		zap.String("order_uuid", event.OrderUuid.String()),
		zap.String("user_uuid", event.UserUuid.String()),
		zap.String("payment_method", event.PaymentMethod),
		zap.String("transaction_uuid", event.TransactionUuid.String()),
	)
	return nil
}
