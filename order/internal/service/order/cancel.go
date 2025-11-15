package order

import (
	"context"

	"github.com/google/uuid"
	"order/internal/model"
)

func (s *service) CancelOrder(ctx context.Context, order_uuid uuid.UUID) error {
	if order_uuid == uuid.Nil {
		return model.ErrEmptyOrderUuid
	}

	reqGetCtx, cancelGet := context.WithTimeout(ctx, model.RequestTimeOutRead)
	defer cancelGet()

	order, err := s.OrderRepository.GetOrder(reqGetCtx, order_uuid)
	if err != nil {
		return err
	}

	if order.OrderStatus == model.OrderStatusPAID {
		return model.ErrCancelOrderStatusPaid
	}

	order.OrderStatus = model.OrderStatusCANCELLED

	reqUpdateCtx, cancelUpdate := context.WithTimeout(ctx, model.RequestTimeOutUpdate)
	defer cancelUpdate()

	return s.OrderRepository.UpdateOrder(reqUpdateCtx, order)
}
