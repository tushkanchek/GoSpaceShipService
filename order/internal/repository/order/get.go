package order

import (
	"context"
	model "order/internal/model"
	"order/internal/repository/converter"
)


func (r *repository) GetOrder(_ context.Context, order_uuid string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exist := r.orders[order_uuid]
	if !exist {
		return nil, model.ErrOrderNotFound
	}

	return converter.OrderToModel(order), nil
}