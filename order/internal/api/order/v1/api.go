package v1

import (
	"order/internal/service"
	orderV1 "shared/pkg/openapi/order/v1"
)



type api struct{
	orderV1.UnimplementedHandler

	service service.OrderService

}

func NewAPI(service service.OrderService) *api{
	return &api{
		service: service,
	}
}