	package order

	import (
		"context"
		"order/internal/model"
	)



	//TODO: check orderstatus cancel
	func (s *service) PayOrder(ctx context.Context, PaymentMethod model.PaymentMethod, order_uuid string) (string, error){
		if order_uuid == ""{
			return "", model.ErrEmptyOrderUuid
		}
		order, err:=s.OrderRepository.GetOrder(ctx, order_uuid)
		if err!=nil{
			return "", err
		}
		if order==nil{
			return "", model.ErrOrderNotFound
		}
		if order.OrderStatus==model.OrderStatusPAID{
			return "", model.ErrPayOrderStatusPaid
		}
		if order.OrderStatus==model.OrderStatusCANCELLED{
			return "", model.ErrPayOrderStatusCancelled
		}
		
		transaction_uuid, err := s.PaymentClient.PayOrder(ctx, order_uuid, order.UserUUID, PaymentMethod)
		if err!=nil{
			return "", err
		}

		order.OrderStatus = model.OrderStatusPAID
		order.TransactionUUID = &transaction_uuid
		order.PaymentMethod = &PaymentMethod

		err = s.OrderRepository.UpdateOrder(ctx, order)
		if err!=nil{
			return "", err
		}

		return transaction_uuid, err

	}
