package converter

import (
	model "order/internal/model"
	repoModel "order/internal/repository/model"
)

func OrderToRepoOrder(order *model.Order) *repoModel.Order {
	return &repoModel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   (*repoModel.PaymentMethod)(order.PaymentMethod),
		OrderStatus:     repoModel.OrderStatus(order.OrderStatus),
	}
}

func OrderToModel(order *repoModel.Order) *model.Order {
	return &model.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   (*model.PaymentMethod)(order.PaymentMethod),
		OrderStatus:     model.OrderStatus(order.OrderStatus),
	}
}
