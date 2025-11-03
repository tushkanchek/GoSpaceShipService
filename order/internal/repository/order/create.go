package order

import (
	"context"
	repoConverter "order/internal/repository/converter"
	model "order/internal/model"
)


func (r *repository) CreateOrder(_ context.Context, order *model.Order) error{
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exist := r.orders[order.OrderUUID]
	if exist{
		return model.ErrOrderAlreadyExists
	}

	r.orders[order.OrderUUID] = repoConverter.OrderToRepoOrder(order)
	return nil
}