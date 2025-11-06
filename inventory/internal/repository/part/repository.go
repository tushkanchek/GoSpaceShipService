package part

import (
	"sync"

	def "inventory/internal/repository"
	repoModel "inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]*repoModel.Part
}

func NewRepository() *repository {
	r := &repository{
		data: make(map[string]*repoModel.Part),
	}

	r.initParts()

	return r
}
