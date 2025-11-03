package order

import (
	repoModel "order/internal/repository/model"
	"sync"
	def "order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)


type repository struct {
	mu     sync.RWMutex
	orders map[string]*repoModel.Order
}

func NewOrderRepository() *repository {
	return &repository{
		orders: make(map[string]*repoModel.Order),
	}
}