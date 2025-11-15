package grpc

import (
	"context"

	"github.com/google/uuid"
	"order/internal/model"
)

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUuid, userUuid string, paymentMethod model.PaymentMethod) (uuid.UUID, error)
}

type InventoryClient interface {
	ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error)
}
