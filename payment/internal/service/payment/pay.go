package payment

import (
	"context"
	"log"

	"github.com/google/uuid"
	"payment/internal/model"
)

func (s *service) PayOrder(_ context.Context, orderUuid, userUuid, PaymentMethod string) (string, error) {
	if orderUuid == "" {
		return "", model.ErrEmptyOrderUuid
	}

	if userUuid == "" {
		return "", model.ErrEmptyUserUuid
	}

	if PaymentMethod == "" {
		return "", model.ErrEmptyPaymentMethod
	}

	transaction_uuid := uuid.NewString()

	log.Printf("Оплата успешно прошла, transaction_uuid: %s\n", transaction_uuid)

	return transaction_uuid, nil
}
