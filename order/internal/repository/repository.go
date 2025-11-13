package repository

import (
	"context"

	model "order/internal/model"

	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	GetOrder(ctx context.Context, order_uuid uuid.UUID) (*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
}
