package app

import (

	paymentV1API "payment/internal/api/payment/v1"
	"payment/internal/service"
	paymentService "payment/internal/service/payment"
	paymentV1 "shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API paymentV1.PaymentServiceServer

	paymentService service.PaymentService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API() paymentV1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = paymentV1API.NewAPI(d.PartService())
	}

	return d.paymentV1API
}

func (d *diContainer) PartService() service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService()
	}

	return d.paymentService
}
