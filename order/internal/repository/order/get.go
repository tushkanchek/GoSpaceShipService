package order

import (
	"context"
	repoConverter "order/internal/repository/converter"
	"log"

	model "order/internal/model"
	repoModel "order/internal/repository/model"

	"github.com/google/uuid"
)

func (r *repository) GetOrder(ctx context.Context, order_uuid uuid.UUID) (*model.Order, error) {
	const op = "repository.order.GetOrder"
	res := r.db.QueryRow(ctx, 

		`SELECT order_uuid, user_uuid, part_uuids, total_price, transaction_uuid, payment_method, order_status
		FROM orders
		WHERE order_uuid=($1);`,	
	order_uuid,
	)

	order:= &repoModel.Order{}
	err := res.Scan(&order.OrderUUID, &order.UserUUID, &order.PartUuids, &order.TotalPrice, &order.TransactionUUID, &order.PaymentMethod, &order.OrderStatus)
	if err!=nil{
		log.Printf("failed to select order: %s: %v\n", op, err)
		return nil, err
	}
	return repoConverter.OrderToModel(order), nil
}
