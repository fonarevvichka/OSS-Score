package util

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetScore(mongoClient *mongo.Client, catalog string, owner string, name string) {
	collection := mongoClient.Database("OSS-Score").Collection(catalog) // TODO MAKE DB NAME ENV VAR

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"owner", owner}},
				bson.D{{"name", name}},
			}},
	}

	res := collection.FindOne(ctx, filter)

	if res.Err() == mongo.ErrNoDocuments { // No match in DB
		fmt.Println("need to do full query")
	} else { // Match in DB found
		var repoInfo RepoInfo
		err := res.Decode(&repoInfo)

		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Doing limited query")
		// Check date and return data or query with according back stop
	}

}
