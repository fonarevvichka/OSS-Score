package main

import (
	"OSS-Score/util"
	"context"
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

	// fmt.Println("Successfully connected and pinged.")

	// src := oauth2.StaticTokenSource(
	// 	&oauth2.Token{AccessToken: os.Getenv("GIT_PAT")},
	// )
	// http_client := oauth2.NewClient(context.Background(), src)

	repoCatalog := "Github"
	repoOwner := "swagger-api"
	repoName := "swagger-ui"

	util.GetScore(mongoClient, repoCatalog, repoOwner, repoName)

	// repoInfo := util.RepoInfo{
	// 	Catalog:      "github",
	// 	Owner:        repoOwner,
	// 	Name:         repoName,
	// 	UpdatedAt:    time.Now(),
	// 	Dependencies: make([]util.Dependency, 0),
	// 	Issues: util.Issues{
	// 		OpenIssues:   make([]util.OpenIssue, 0),
	// 		ClosedIssues: make([]util.ClosedIssue, 0),
	// 	},
	// }

	// err = util.GetCoreRepoInfo(http_client, &repoInfo)

	// startDate := time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC).Format(time.RFC3339)
	// util.GetGithubIssues(http_client, &repoInfo, startDate)
	// util.GetGithubDependencies(http_client, &repoInfo)

	// collection := mongoClient.Database("OSS-Score").Collection("Github")

	// result, err := collection.InsertOne(context.TODO(), repoInfo)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(result)
}
