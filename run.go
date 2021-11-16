package main

import (
	db "OSS-Score/mongo"
	"OSS-Score/util"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/oauth2"
)

func main() {

	uri := os.Getenv("MONGO_URI")
	// Create a new mongo_client and connect to the server
	mongo_client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	defer func() {
		if err = mongo_client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Ping the primary
	if err := mongo_client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_PAT")},
	)
	http_client := oauth2.NewClient(context.Background(), src)

	repoOwner := "swagger-api"
	repoName := "swagger-ui"

	repoInfo := util.RepoInfo{
		Catalog:      "github",
		Owner:        repoOwner,
		Name:         repoName,
		Dependencies: make([]util.Dependency, 0),
		Issues: util.Issues{
			OpenIssues:   make([]util.OpenIssue, 0),
			ClosedIssues: make([]util.ClosedIssue, 0),
		},
	}

	err = util.GetCoreRepoInfo(http_client, &repoInfo)

	startDate := time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC).Format(time.RFC3339)
	util.GetGithubIssues(http_client, &repoInfo, startDate)
	util.GetGithubDependencies(http_client, &repoInfo)

	collection := mongo_client.Database("OSS-Score").Collection("Github")
	result, err := db.InsertNewRepo(*collection, context.TODO(), repoInfo)

	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(result)
}
