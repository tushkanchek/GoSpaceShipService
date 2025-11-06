package part

import (
	"context"

	"inventory/internal/model"
	repoConverter "inventory/internal/repository/converter"
)

func (r *repository) GetPart(_ context.Context, partUuid string) (*model.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	part, ok := r.data[partUuid]
	if !ok {
		return nil, model.ErrPartNotFound
	}

	return repoConverter.RepoPartToModel(*part), nil
}
