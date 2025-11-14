package order

import (
	"context"

	"github.com/google/uuid"
	"order/internal/model"
)

func (s *service) GetOrderByUUID(ctx context.Context, order_uuid uuid.UUID) (*model.Order, error) {
	if order_uuid == uuid.Nil {
		return nil, model.ErrEmptyOrderUuid
	}

	order, err := s.OrderRepository.GetOrder(ctx, order_uuid)
	if err != nil {
		return nil, err
	}

	return order, nil
}
