package grpc

import (
	"context"

	"order/internal/model"
)

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUuid, userUuid string, paymentMethod model.PaymentMethod) (string, error)
}

type InventoryClient interface {
	ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error)
}
