package service

import (
	"context"

	"order/internal/model"
)

type OrderService interface {
	GetOrderByUUID(ctx context.Context, order_uuid string) (*model.Order, error)
	CreateOrder(ctx context.Context, user_uuid string, part_uuids []string) (*model.Order, error)
	PayOrder(ctx context.Context, PaymentMethod model.PaymentMethod, order_uuid string) (string, error)
	CancelOrder(ctx context.Context, order_uuid string) error
}
