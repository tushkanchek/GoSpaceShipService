package service

import "context"




type PaymentService interface{
	PayOrder(ctx context.Context, orderUuid string, userUuid string, PaymentMethod string) (string, error)
}
