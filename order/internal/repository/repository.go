package repository

import (
	"context"

	model "order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	GetOrder(ctx context.Context, order_uuid string) (*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
}
