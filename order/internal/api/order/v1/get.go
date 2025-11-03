package v1

import (
	"context"
	"errors"
	"net/http"

	"order/internal/converter"
	"order/internal/model"
	orderV1 "shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order, err := a.service.GetOrderByUUID(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrEmptyOrderUuid) {
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "Order uuid is empty",
			}, nil
		}
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Order for uuid '" + params.OrderUUID + "' not found",
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "get order: Unexpected error: " + err.Error(),
		}, nil

	}

	return &orderV1.GetOrderResponse{
		Order: *converter.OrderToApi(order),
	}, nil
}
