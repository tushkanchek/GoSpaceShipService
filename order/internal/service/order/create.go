package order

import (
	"context"

	"github.com/google/uuid"
	model "order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, user_uuid uuid.UUID, part_uuids []uuid.UUID) (*model.Order, error) {
	if user_uuid == uuid.Nil {
		return nil, model.ErrEmptyUserUuid
	}
	if len(part_uuids) == 0 {
		return nil, model.ErrEmptyListPartUuids
	}

	invPartUuids := make([]string, 0, len(part_uuids))
	for _, elem := range part_uuids {
		invPartUuids = append(invPartUuids, elem.String())
	}

	reqListCtx, cancelList := context.WithTimeout(ctx, model.RequestTimeOutRead)
	defer cancelList()

	resp, err := s.InventoryClient.ListParts(reqListCtx, &model.PartsFilter{
		Uuids: invPartUuids,
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
		OrderUUID:   uuid.MustParse(uuid.NewString()),
		UserUUID:    user_uuid,
		PartUuids:   part_uuids,
		TotalPrice:  totalPrice,
		OrderStatus: model.OrderStatusPENDINGPAYMENT,
	}

	reqCreateCtx, cancelCreate := context.WithTimeout(ctx, model.RequestTimeOutUpdate)
	defer cancelCreate()

	err = s.OrderRepository.CreateOrder(reqCreateCtx, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}
