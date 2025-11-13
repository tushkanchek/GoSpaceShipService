package v1

import (
	"context"
	"net/http"

	"order/internal/model"
	orderV1 "shared/pkg/openapi/order/v1"

	"github.com/google/uuid"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order_uuid, err := uuid.Parse(params.OrderUUID)
	if err!=nil{
		return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "order uuid couldnt be parsed",
			}, nil
	}
	err = a.service.CancelOrder(ctx, order_uuid)
	if err != nil {
		switch err {
		case model.ErrOrderNotFound:
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		case model.ErrEmptyOrderUuid:
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		case model.ErrCancelOrderStatusPaid:
			return &orderV1.ConflictError{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}, nil
		default:
			return &orderV1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}, nil
		}
	}
	return &orderV1.CancelOrderNoContent{}, nil
}
