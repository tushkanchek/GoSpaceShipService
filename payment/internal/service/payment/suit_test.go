package payment

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"platform/pkg/logger"
)

type ServiceSuite struct {
	suite.Suite

	service *service
}

func (s *ServiceSuite) SetupTest() {
	logger.SetNopLogger()

	s.service = NewService()
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
