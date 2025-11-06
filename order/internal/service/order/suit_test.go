package order

import (
	"context"
	repoMocks "order/internal/repository/mocks"
	"testing"

	"order/internal/client/grpc/mocks"

	"github.com/stretchr/testify/suite"
)




type ServiceSuite struct{
	suite.Suite

	ctx context.Context

	orderRepo *repoMocks.OrderRepository
	inventoryClient *mocks.InventoryClient
	paymentClient *mocks.PaymentClient

	service *service
}

func (s *ServiceSuite) SetupTest(){
	s.ctx = context.Background()

	s.orderRepo = repoMocks.NewOrderRepository(s.T())
	s.inventoryClient = mocks.NewInventoryClient(s.T())
	s.paymentClient = mocks.NewPaymentClient(s.T())	

	s.service = NewService(
		s.orderRepo,
		s.inventoryClient, 
		s.paymentClient,
	)
}

func (s *ServiceSuite) TearDownTest(){
}

func TestServiceIntegration(t *testing.T){
	suite.Run(t, new(ServiceSuite))
}