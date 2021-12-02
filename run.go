package main

import (
	"OSS-Score/util"
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {

	uri := os.Getenv("MONGO_URI")
	// Create a new mongo_client and connect to the server
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Ping the primary
	if err := mongoClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")

	repoCatalog := "github"
	// repoOwner := "swagger-api"
	// repoName := "swagger-ui"
	repoOwner := "facebook"
	repoName := "react"
	// repoOwner := "jasonlong"
	// repoName := "isometric-contributions"
	// repoOwner := "fonarevvichka"
	// repoName := "OSS-Score"

	repoInfoResponse := util.GetScore(mongoClient, repoCatalog, repoOwner, repoName, 1)
	fmt.Println(repoInfoResponse.Ready)
}
