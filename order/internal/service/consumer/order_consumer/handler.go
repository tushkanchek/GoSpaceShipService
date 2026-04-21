package orderconsumer

import (
	"context"
	"order/internal/model"
	"platform/pkg/kafka/consumer"
	"platform/pkg/logger"

	"go.uber.org/zap"
)




func (s *service) OrderHandler(ctx context.Context, msg consumer.Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	event, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderAssembled", zap.Error(err))
		return err
	}
	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUuid.String()),
		zap.String("order_uuid", event.OrderUuid.String()),
		zap.String("user_uuid", event.UserUuid.String()),
		zap.Int32("build_time_sec", event.BuildTimeSec),
	)

	logger.Info(ctx, "Update order status",
		zap.String("order_uuid", event.OrderUuid.String()),
	)

	err =s.orderService.UpdateOrderStatus(ctx, event.OrderUuid, model.OrderStatusASSEMBLED)
	if err != nil {
		logger.Error(ctx, "Failed to update order status", zap.Error(err))
		return err
	}

	return nil
}