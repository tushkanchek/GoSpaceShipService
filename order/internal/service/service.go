package service

import (
	"context"

	"order/internal/model"
)

type OrderService interface {
	GetOrderByUUID(ctx context.Context, order_uuid string) (*model.Order, error)
	CreateOrder(ctx context.Context, user_uuid string, part_uuids []string) (*model.Order, error)
	PayOrder(ctx context.Context, order_uuid string, PaymentMethod model.PaymentMethod) (string, error)
	CancelOrder(ctx context.Context, order_uuid string) error
}
