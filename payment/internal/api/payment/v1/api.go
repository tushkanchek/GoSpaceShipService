package v1

import (
	service "payment/internal/service"
	paymentV1 "shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

func NewAPI(paymentServce service.PaymentService) *api {
	return &api{
		paymentService: paymentServce,
	}
}
