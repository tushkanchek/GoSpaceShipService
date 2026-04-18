package order

import (
	"context"

	"github.com/google/uuid"
	model "order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, userUUID uuid.UUID, partUUIDs []uuid.UUID) (*model.Order, error) {
	if userUUID == uuid.Nil {
		return nil, model.ErrEmptyUserUuid
	}
	if len(partUUIDs) == 0 {
		return nil, model.ErrEmptyListPartUuids
	}

	invPartUuids := make([]string, 0, len(partUUIDs))
	for _, elem := range partUUIDs {
		invPartUuids = append(invPartUuids, elem.String())
	}

	reqListCtx, cancelList := context.WithTimeout(ctx, model.RequestTimeOutRead)
	defer cancelList()

	resp, err := s.InventoryClient.ListParts(reqListCtx, &model.PartsFilter{
		Uuids: invPartUuids,
	})
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, model.ErrPartsByUuidsNotFound
	}

	var totalPrice float64
	for _, p := range resp {
		totalPrice += p.Price
	}

	order := &model.Order{
		OrderUUID:   uuid.MustParse(uuid.NewString()),
		UserUUID:    userUUID,
		PartUuids:   partUUIDs,
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
