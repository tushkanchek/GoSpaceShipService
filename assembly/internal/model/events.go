package model

import "github.com/google/uuid"

type OrderPaid struct{
	EventUuid uuid.UUID
	OrderUuid uuid.UUID
	UserUuid uuid.UUID
	PaymentMethod string
	transaction_uuid uuid.UUID
}
