package order

import (
	"context"

	"github.com/google/uuid"
	model "order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, user_uuid string, part_uuids []string) (*model.Order, error) {
	if user_uuid == "" {
		return nil, model.ErrEmptyUserUuid
	}
	if len(part_uuids) == 0 {
		return nil, model.ErrEmptyListPartUuids
	}

	resp, err := s.InventoryClient.ListParts(ctx, &model.PartsFilter{
		Uuids: part_uuids,
	})
	if len(resp) == 0 {
		return nil, model.ErrPartsByUuidsNotFound
	}
	if err != nil {
		return nil, err
	}

	var totalPrice float64 = 0
	for _, p := range resp {
		totalPrice += p.Price
	}

	order := &model.Order{
		OrderUUID:   uuid.NewString(),
		UserUUID:    user_uuid,
		PartUuids:   part_uuids,
		TotalPrice:  totalPrice,
		OrderStatus: model.OrderStatusPENDINGPAYMENT,
	}

	err = s.OrderRepository.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}
