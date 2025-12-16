package order

import (
	"testing"

	"github.com/stretchr/testify/suite"
	grpcMocks "order/internal/client/grpc/mocks"
	repoMocks "order/internal/repository/mocks"
	"order/internal/service/mocks"
)

type ServiceSuite struct {
	suite.Suite

	orderRepo         *repoMocks.OrderRepository
	inventoryClient   *grpcMocks.InventoryClient
	paymentClient     *grpcMocks.PaymentClient
	orderPaidProducer *mocks.OrderProducerService

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.orderRepo = repoMocks.NewOrderRepository(s.T())
	s.inventoryClient = grpcMocks.NewInventoryClient(s.T())
	s.paymentClient = grpcMocks.NewPaymentClient(s.T())
	s.orderPaidProducer = mocks.NewOrderProducerService(s.T())

	s.service = NewService(
		s.orderRepo,
		s.inventoryClient,
		s.paymentClient,
		s.orderPaidProducer,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
