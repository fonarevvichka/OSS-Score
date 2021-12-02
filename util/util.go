package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
)

func GetScore(mongoClient *mongo.Client, catalog string, owner string, name string, level int) RepoInfoDBResponse {
	collection := mongoClient.Database("OSS-Score").Collection(catalog) // TODO MAKE DB NAME ENV VAR
	shelfLife := 7                                                      // Days TODO: make env var

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
	var repoInfo RepoInfo
	infoReady := false

	if res.Err() == mongo.ErrNoDocuments { // No match in DB
		fmt.Println("need to do full query")
		repoInfo = queryGithub(catalog, owner, name, time.Now().AddDate(-1, 0, 0)) // hardcode to 1 year timefirame
		infoReady = true                                                           // temp while this is synchronous
		_, err := collection.InsertOne(context.TODO(), repoInfo)
		fmt.Println("Inserting Data")
		if err != nil {
			log.Fatal(err)
		}

	} else { // Match in DB found
		err := res.Decode(&repoInfo)

		if err != nil {
			log.Fatalln(err)
		}

		expireDate := time.Now().AddDate(0, 0, -shelfLife)
		if repoInfo.UpdatedAt.Before(expireDate) {
			fmt.Println("out of date: need to make partial query")

			repoInfo := queryGithub(catalog, owner, name, repoInfo.UpdatedAt) // pull 1 week of data
			infoReady = true                                                  // temp while this is synchronous

			insertableData := bson.D{primitive.E{Key: "$set", Value: repoInfo}}
			_, err := collection.UpdateOne(context.TODO(), filter, insertableData)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			infoReady = true
		}
	}

	if level < 1 {
		level += 1
		fmt.Println(level)
		for _, dependency := range repoInfo.Dependencies {
			fmt.Println("updating: " + dependency.Name)
			GetScore(mongoClient, dependency.Catalog, dependency.Owner, dependency.Name, level)
		}
	}

	res = collection.FindOne(ctx, filter)
	activityScore, activityConfidence := CalculateActivityScore(&repoInfo, time.Now().AddDate(-1, 0, 0)) // startpoint hardcoded for now
	repoInfo.ActivityScore = float32(activityScore)
	repoInfo.ActivityConfidence = float32(activityConfidence)
	insertableData := bson.D{primitive.E{Key: "$set", Value: repoInfo}}
	_, err := collection.UpdateOne(context.TODO(), filter, insertableData)
	if err != nil {
		log.Fatal(err)
	}

	return RepoInfoDBResponse{
		Ready:    infoReady,
		RepoInfo: RepoInfo{},
	}
}

func queryGithub(catalog string, owner string, name string, startPoint time.Time) RepoInfo {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_PAT")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	repoInfo := RepoInfo{
		Catalog:      catalog,
		Owner:        owner,
		Name:         name,
		UpdatedAt:    time.Now(),
		Dependencies: make([]Dependency, 0),
		Issues: Issues{
			OpenIssues:   make([]OpenIssue, 0),
			ClosedIssues: make([]ClosedIssue, 0),
		},
	}

	// These need to be async
	GetCoreRepoInfo(httpClient, &repoInfo)
	GetGithubIssues(httpClient, &repoInfo, startPoint.Format(time.RFC3339))
	GetGithubDependencies(httpClient, &repoInfo)
	GetGithubCommits(httpClient, &repoInfo, startPoint.Format(time.RFC3339))
	GetGithubReleases(httpClient, &repoInfo, startPoint.Format(time.RFC3339))

	return repoInfo
}

func dependencyInSlice(dependency Dependency, dependencies []Dependency) bool {
	for _, elem := range dependencies {
		if dependency.Catalog == elem.Catalog &&
			dependency.Owner == elem.Owner &&
			dependency.Name == elem.Name &&
			dependency.Version == elem.Version { // NOT SURE IF DEEP COMPARE LIKE THIS IS NEEDED
			return true
		}
	}
	return false
}
