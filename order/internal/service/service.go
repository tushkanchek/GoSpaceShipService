package service

import (
	"context"

	"order/internal/model"

	"github.com/google/uuid"
)

type OrderService interface {
	GetOrderByUUID(ctx context.Context, order_uuid uuid.UUID) (*model.Order, error)
	CreateOrder(ctx context.Context, user_uuid uuid.UUID, part_uuids []uuid.UUID) (*model.Order, error)
	PayOrder(ctx context.Context, order_uuid uuid.UUID, PaymentMethod model.PaymentMethod) (uuid.UUID, error)
	CancelOrder(ctx context.Context, order_uuid uuid.UUID) error
}
