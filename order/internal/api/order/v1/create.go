package v1

import (
	"context"
	"net/http"
	"order/internal/model"
	orderV1 "shared/pkg/openapi/order/v1"
)



func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error){
	order, err := a.service.CreateOrder(ctx, req.UserUUID, req.PartUuids)
	if err!=nil{
		switch err{
		case model.ErrOrderAlreadyExists, model.ErrEmptyUserUuid, model.ErrEmptyListPartUuids:
			return &orderV1.BadRequestError{  //TODO: return make it conflict error
				Code: http.StatusBadRequest,
				Message: err.Error(),
			}, nil
		case model.ErrPartsByUuidsNotFound:
			return &orderV1.NotFoundError{
				Code: http.StatusNotFound,
				Message: err.Error(),
		}, nil	
		}
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID: order.OrderUUID,
		TotalPrice: float32(order.TotalPrice),
	}, nil

}

//func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
// 	resp, err := h.inventoryClient.ListParts(
// 		ctx,
// 		&inventoryV1.ListPartsRequest{
// 			Filter: &inventoryV1.PartsFilter{
// 				Uuids: req.PartUuids,
// 			},
// 		},
// 	)
// 	if err != nil {
// 		return &orderV1.InternalServerError{
// 			Code:    500,
// 			Message: fmt.Sprintf("failed to get parts info: %v", err),
// 		}, nil
// 	}
// 	if resp == nil {
// 		return &orderV1.NotFoundError{
// 			Code:    404,
// 			Message: "parts not found",
// 		}, nil
// 	}
// 	if len(resp.Parts) < len(req.PartUuids) {
// 		return &orderV1.BadRequestError{
// 			Code:    400,
// 			Message: "some parts not found",
// 		}, nil
// 	}

// 	var totalPrice float64
// 	for _, p := range resp.Parts {
// 		totalPrice += p.Price
// 	}

// 	order := &orderV1.Order{
// 		OrderUUID:  uuid.NewString(),
// 		UserUUID:   req.GetUserUUID(),
// 		PartUuids:  req.GetPartUuids(),
// 		TotalPrice: float32(totalPrice),
// 		Status:     orderV1.OrderStatusPENDINGPAYMENT,
// 	}

// 	h.storage.CreateOrder(order)
// 	return &orderV1.CreateOrderResponse{
// 		OrderUUID:  order.OrderUUID,
// 		TotalPrice: order.TotalPrice,
// 	}, nil
// }