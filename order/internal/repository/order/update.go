package order

import (
	"context"

	model "order/internal/model"
	repoConverter "order/internal/repository/converter"
)

func (r *repository) UpdateOrder(_ context.Context, order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exist := r.orders[order.OrderUUID]
	if !exist {
		return model.ErrOrderNotFound
	}

	r.orders[order.OrderUUID] = repoConverter.OrderToRepoOrder(order)
	return nil
}
