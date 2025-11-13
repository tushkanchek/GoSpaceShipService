package order

import (
	"context"

	"order/internal/model"

	"github.com/google/uuid"
)

func (s *service) CancelOrder(ctx context.Context, order_uuid uuid.UUID) error {
	if order_uuid == uuid.Nil {
		return model.ErrEmptyOrderUuid
	}

	order, err := s.OrderRepository.GetOrder(ctx, order_uuid)
	if err != nil {
		return err
	}

	if order.OrderStatus == model.OrderStatusPAID {
		return model.ErrCancelOrderStatusPaid
	}

	order.OrderStatus = model.OrderStatusCANCELLED
	return s.OrderRepository.UpdateOrder(ctx, order)
}
