package service

import "context"

type PaymentService interface {
	PayOrder(ctx context.Context, orderUuid, userUuid, PaymentMethod string) (string, error)
}
