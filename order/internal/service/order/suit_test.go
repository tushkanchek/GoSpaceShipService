package order

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"order/internal/client/grpc/mocks"
	repoMocks "order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	orderRepo       *repoMocks.OrderRepository
	inventoryClient *mocks.InventoryClient
	paymentClient   *mocks.PaymentClient

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.orderRepo = repoMocks.NewOrderRepository(s.T())
	s.inventoryClient = mocks.NewInventoryClient(s.T())
	s.paymentClient = mocks.NewPaymentClient(s.T())

	s.service = NewService(
		s.orderRepo,
		s.inventoryClient,
		s.paymentClient,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
