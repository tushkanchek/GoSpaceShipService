package order

import (
	"order/internal/client/grpc"
	"order/internal/repository"
	def "order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	OrderRepository repository.OrderRepository
	InventoryClient grpc.InventoryClient
	PaymentClient   grpc.PaymentClient
	orderProducerService def.OrderProducerService
}




func NewService(orderRepository repository.OrderRepository, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient, orderProducerService def.OrderProducerService) *service {
	return &service{
		OrderRepository: orderRepository,
		InventoryClient: inventoryClient,
		PaymentClient:   paymentClient,
		orderProducerService: orderProducerService,
	}
}
