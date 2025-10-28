package model

import "errors"


var (
	ErrEmptyOrderUuid = errors.New("order uuid is empty")

	ErrEmptyUserUuid = errors.New("user uuid is empty")

	ErrEmptyPaymentMethod = errors.New("payment method is empty")
)