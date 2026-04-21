package order

import (
	"context"
	"platform/pkg/logger"

	"order/internal/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *service) UpdateOrderStatus(ctx context.Context, order_uuid uuid.UUID, status model.OrderStatus) error {
	order, err := s.OrderRepository.GetOrder(ctx, order_uuid)
	if err != nil {
		return err
	}

	order.OrderStatus = status

	err = s.OrderRepository.UpdateOrder(ctx, order)
	if err != nil {
		logger.Error(ctx, "Failed to update order in repository", zap.String("order_uuid", order_uuid.String()), zap.Error(err))
		return err
	}

	return nil
}		