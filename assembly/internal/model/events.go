package model

import "github.com/google/uuid"

type OrderPaidEvent struct {
	EventUuid       uuid.UUID
	OrderUuid       uuid.UUID
	UserUuid        uuid.UUID
	PaymentMethod   string
	TransactionUuid uuid.UUID
}
