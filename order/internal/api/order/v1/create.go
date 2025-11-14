package v1

import (
	"context"
	"errors"
	"net/http"

	"order/internal/model"
	orderV1 "shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	order, err := a.service.CreateOrder(ctx, req.UserUUID, req.PartUuids)
	if err != nil {
		if errors.Is(err, model.ErrOrderAlreadyExists) && errors.Is(err, model.ErrEmptyUserUuid) && errors.Is(err, model.ErrEmptyListPartUuids) {
			return &orderV1.BadRequestError{ // TODO: return make it conflict error
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		}
		if errors.Is(err, model.ErrPartsByUuidsNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		}
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: float32(order.TotalPrice),
	}, nil
}
