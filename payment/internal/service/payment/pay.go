package payment

import (
	"context"
	"log"
	"payment/internal/model"

	"github.com/google/uuid"
)


func (s *service) PayOrder(_ context.Context, orderUuid string, userUuid string, PaymentMethod string) (string, error){
	if orderUuid == "" {
		return "", model.ErrEmptyOrderUuid
	}

	if userUuid == ""{
		return "", model.ErrEmptyUserUuid
	}

	if PaymentMethod == "" {
		return "", model.ErrEmptyPaymentMethod
	}

	transaction_uuid := uuid.NewString()

	log.Printf("Оплата успешно прошла, transaction_uuid: %s\n", transaction_uuid)

	return transaction_uuid, nil
}