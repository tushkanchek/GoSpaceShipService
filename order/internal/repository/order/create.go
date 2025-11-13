package order

import (
	"context"
	"log"

	model "order/internal/model"
	repoConverter "order/internal/repository/converter"
)

func (r *repository) CreateOrder(ctx context.Context, order *model.Order) error {
	const op = "repository.order.CreateOrder"
	repoOrder := repoConverter.OrderToRepoOrder(order)

	_, err := r.db.Exec(ctx, 
	`
		INSERT INTO orders (order_uuid, user_uuid, part_uuids, total_price, transaction_uuid, payment_method, order_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`, repoOrder.OrderUUID, repoOrder.UserUUID, repoOrder.PartUuids, repoOrder.TotalPrice, repoOrder.TransactionUUID, repoOrder.PaymentMethod, repoOrder.OrderStatus,
	)
	if err!=nil{
		log.Printf("Failed to insert values to orders: %s %v\n", op, err)
		return err
	}
	return nil
}
