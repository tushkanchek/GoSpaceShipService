package part

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	inventoryRepo *mocks.InventoryRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.inventoryRepo = mocks.NewInventoryRepository(s.T())

	s.service = NewService(
		s.inventoryRepo,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
