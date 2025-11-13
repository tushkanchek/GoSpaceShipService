package v1

import (
	"context"
	"net/http"

	"order/internal/model"
	orderV1 "shared/pkg/openapi/order/v1"

	"github.com/google/uuid"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	user_uuid, err := uuid.Parse(req.UserUUID)
	if err!=nil{
		return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "User uuid couldn't be parsed",
			}, nil
	}

	part_uuids := make([]uuid.UUID, 0, len(req.PartUuids))
	for _, el := range req.PartUuids{
		part_uuid, err := uuid.Parse(el)
		if err!=nil{
			return &orderV1.BadRequestError{
					Code:    http.StatusBadRequest,
					Message: "parts uuid couldnt be parsed",
				}, nil
		}
		part_uuids = append(part_uuids, part_uuid)
	}
	order, err := a.service.CreateOrder(ctx, user_uuid, part_uuids)
	if err != nil {
		switch err {
		case model.ErrOrderAlreadyExists, model.ErrEmptyUserUuid, model.ErrEmptyListPartUuids:
			return &orderV1.BadRequestError{ // TODO: return make it conflict error
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		case model.ErrPartsByUuidsNotFound:
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		}
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID.String(),
		TotalPrice: float32(order.TotalPrice),
	}, nil
}

