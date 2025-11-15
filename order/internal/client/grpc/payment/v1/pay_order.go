package v1

import (
	"context"

	"github.com/google/uuid"
	"order/internal/model"
	paymentV1 "shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUuid, userUuid string, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	// TODO: add workd with ctx

	resp, err := c.generatedClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     orderUuid,
		UserUuid:      userUuid,
		PaymentMethod: paymentV1.PaymentMethod(paymentMethod),
	})
	if err != nil {
		return uuid.Nil, err
	}

	transaction_uuid, err := uuid.Parse(resp.TransactionUuid)
	if err != nil {
		return uuid.Nil, err
	}
	return transaction_uuid, nil
}
