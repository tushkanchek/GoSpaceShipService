package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"payment/internal/model"
	paymentV1 "shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transaction_uuid, err := a.paymentService.PayOrder(ctx, req.OrderUuid, req.UserUuid, req.PaymentMethod.String())
	if err != nil {
		if errors.Is(err, model.ErrEmptyOrderUuid) {
			return nil, status.Errorf(codes.InvalidArgument, "order uuid is empty")
		}
		if errors.Is(err, model.ErrEmptyUserUuid) {
			return nil, status.Errorf(codes.InvalidArgument, "user uuid is empty")
		}
		if errors.Is(err, model.ErrEmptyPaymentMethod) {
			return nil, status.Errorf(codes.InvalidArgument, "payment method uuid is empty")
		}

		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}
	return &paymentV1.PayOrderResponse{
		TransactionUuid: transaction_uuid,
	}, nil
}
