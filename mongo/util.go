package mongo

import (
	"context"
	"log"

	"go_exploring/util"

	"go.mongodb.org/mongo-driver/mongo"
)

func InsertNewRepo(collection mongo.Collection, ctx context.Context, repo util.RepoInfo) (*mongo.InsertOneResult, error) {
	insertResult, err := collection.InsertOne(ctx, repo)

	if err != nil {
		log.Fatal(err)
	}

	return insertResult, err
}
