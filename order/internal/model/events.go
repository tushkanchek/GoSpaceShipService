package model

import "github.com/google/uuid"

type OrderPaidEvent struct {
	EventUuid       uuid.UUID
	OrderUuid       uuid.UUID
	UserUuid        uuid.UUID
	PaymentMethod   int32
	TransactionUuid uuid.UUID
}

type OrderAssembledEvent struct {
	EventUuid 		uuid.UUID
	OrderUuid 		uuid.UUID
	UserUuid  		uuid.UUID
	BuildTimeSec 	int32
}