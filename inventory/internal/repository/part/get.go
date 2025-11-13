package part

import (
	"context"
	"errors"
	model "inventory/internal/model"
	repoConverter "inventory/internal/repository/converter"
	repoModel "inventory/internal/repository/model"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetPart(ctx context.Context, partUuid string) (*model.Part, error) {
	const op = "inventory.internal.repository.GetPart"
	var part *repoModel.Part
	err := r.collection.FindOne(ctx, bson.M{"uuid": partUuid}).Decode(&part)
	if err!=nil{
		if errors.Is(err, mongo.ErrNoDocuments){
			log.Printf("not found part %s %s:%v", partUuid, op, err)
			return nil, model.ErrPartNotFound
		}
		log.Printf("failed to get part %s with error %s: %v", partUuid, op, err)
		return nil, err
	}
	return repoConverter.RepoPartToModel(*part), nil
}
