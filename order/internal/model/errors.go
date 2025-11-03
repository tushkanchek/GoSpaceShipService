package model

import "errors"


var (
	ErrOrderAlreadyExists = errors.New("order with this uuid already exists")

	ErrOrderNotFound = errors.New("order not found")

	ErrEmptyOrderUuid = errors.New("order uuid is empty")

	ErrEmptyUserUuid = errors.New("user uuid is empty")

	ErrPartsByUuidsNotFound = errors.New("all parts with these uuids were not found")

	ErrCancelOrderStatusPaid = errors.New("can't cancel already paid order")

	ErrPayOrderStatusPaid = errors.New("can't pay already paid order")

	ErrPayOrderStatusCancelled = errors.New("can't pay already cancelled order")

	ErrEmptyListPartUuids = errors.New("list of part uuids is empty")
)