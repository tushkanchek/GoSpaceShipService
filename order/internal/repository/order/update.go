package order

import (
	"context"
	"log"

	model "order/internal/model"
	repoConverter "order/internal/repository/converter"
)

func (r *repository) UpdateOrder(ctx context.Context, order *model.Order) error {
	const op = "repository.Order.UpdateOrder"

	repoOrder := repoConverter.OrderToRepoOrder(order)

	_, err := r.GetOrder(ctx, repoOrder.OrderUUID)
	if err != nil {
		log.Printf("Order not exist yet: %s: %v\n", op, err)
		return err
	}

	_, err = r.db.Exec(ctx,
		`UPDATE orders
		SET part_uuids = ($1), 
		total_price = ($2),
		transaction_uuid = ($3),
		payment_method = ($4),
		order_status = ($5)
		WHERE order_uuid = ($6);
		`, repoOrder.PartUuids, repoOrder.TotalPrice, repoOrder.TransactionUUID, repoOrder.PaymentMethod, repoOrder.OrderStatus, repoOrder.OrderUUID,
	)
	if err != nil {
		log.Printf("Failed to update order: %s %v\n", op, err)
		return err
	}
	log.Printf("Successfully update order with UUID: %s\n", repoOrder.OrderUUID.String())
	return nil
}
