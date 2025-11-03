package order

import (
	"context"

	"order/internal/model"
)

func (s *service) GetOrderByUUID(ctx context.Context, order_uuid string) (*model.Order, error) {
	if order_uuid == "" {
		return nil, model.ErrEmptyOrderUuid
	}

	order, err := s.OrderRepository.GetOrder(ctx, order_uuid)
	if err != nil {
		return nil, err
	}

	return order, nil
}
