package v1

import (
	"context"
	"net/http"

	"order/internal/converter"
	"order/internal/model"

	orderV1 "shared/pkg/openapi/order/v1"
)




func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error){
	transaction_uuid, err := a.service.PayOrder(ctx, converter.PaymentMethodApiToModel(req.PaymentMethod), params.OrderUUID)
	if err!=nil{
		switch err{
		case model.ErrOrderNotFound:
			return &orderV1.NotFoundError{
				Code: http.StatusNotFound,
				Message: err.Error(),
			}, nil
		
		case model.ErrPayOrderStatusPaid, model.ErrPayOrderStatusPaid, model.ErrEmptyOrderUuid:
			return &orderV1.BadRequestError{
				Code: http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		default:
			return &orderV1.BadRequestError{
				Code: http.StatusInternalServerError,
				Message: "pay order pay unexpected error: " + err.Error(),
			}, nil
		}
	

	}
	return &orderV1.PayOrderResponse{
		TransactionUUID: transaction_uuid,
	}, nil

}


// func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
// 	order := h.storage.GetOrder(params.OrderUUID)
// 	if order == nil {
// 		return &orderV1.NotFoundError{
// 			Code:    404,
// 			Message: "order with uuid '" + params.OrderUUID + "' not found",
// 		}, nil
// 	}
// 	transaction_uuid, err := h.paymentClient.PayOrder(
// 		ctx,
// 		&paymentV1.PayOrderRequest{
// 			OrderUuid:     order.OrderUUID,
// 			UserUuid:      order.UserUUID,
// 			PaymentMethod: paymentV1.PaymentMethod(req.PaymentMethod),
// 		},
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("payment method error: %v", err)
// 	}