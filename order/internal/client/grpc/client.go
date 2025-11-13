package grpc

import (
	"context"

	"order/internal/model"

	"github.com/google/uuid"
)

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUuid, userUuid string, paymentMethod model.PaymentMethod) (uuid.UUID, error)
}

type InventoryClient interface {
	ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error)
}
