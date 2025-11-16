package payment

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"payment/internal/model"
)

const (
	PayDelay = 1 * time.Second
)

func (s *service) PayOrder(ctx context.Context, orderUuid, userUuid, paymentMethod string) (string, error) {
	if orderUuid == "" {
		return "", model.ErrEmptyOrderUuid
	}

	if userUuid == "" {
		return "", model.ErrEmptyUserUuid
	}

	if paymentMethod == "" {
		return "", model.ErrEmptyPaymentMethod
	}

	if dl, ok := ctx.Deadline(); ok {
		log.Printf("⌛ Context with deadline: %v\n", time.Until(dl))
	} else {
		log.Printf("⌛ Context with no timeout\n")
	}

	timer := time.NewTimer(PayDelay)
	defer timer.Stop()
	select {
	case <-timer.C:
		transaction_uuid := uuid.NewString()

		log.Printf("Оплата успешно прошла, transaction_uuid: %s\n", transaction_uuid)

		return transaction_uuid, nil
	case <-ctx.Done():
		return "", ctx.Err()

	}
}
