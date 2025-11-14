package v1

import (
	"context"
	"net/http"

	"order/internal/converter"
	"order/internal/model"
	orderV1 "shared/pkg/openapi/order/v1"

	
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	transaction_uuid, err := a.service.PayOrder(ctx, params.OrderUUID, converter.PaymentMethodApiToModel(req.PaymentMethod))
	if err != nil {
		switch err {
		case model.ErrOrderNotFound:
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil

		case model.ErrPayOrderStatusPaid, model.ErrPayOrderStatusPaid, model.ErrEmptyOrderUuid:
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		default:
			return &orderV1.BadRequestError{
				Code:    http.StatusInternalServerError,
				Message: "pay order pay unexpected error: " + err.Error(),
			}, nil
		}
	}
	return &orderV1.PayOrderResponse{
		TransactionUUID: transaction_uuid,
	}, nil
}

