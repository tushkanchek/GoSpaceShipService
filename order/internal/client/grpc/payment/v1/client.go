package v1

import (
	"time"

	def "order/internal/client/grpc"
	paymentV1 "shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	timeout         time.Duration
	generatedClient paymentV1.PaymentServiceClient
}

func NewClient(timeout time.Duration, generatedClient paymentV1.PaymentServiceClient) *client {
	return &client{
		timeout:         timeout,
		generatedClient: generatedClient,
	}
}
