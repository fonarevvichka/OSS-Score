package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/oauth2"
)

func GetRepoFromDB(collection *mongo.Collection, owner string, name string) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"owner", owner}},
				bson.D{{"name", name}},
			}},
	}

	return collection.FindOne(ctx, filter)
}

func GetScore(mongoClient *mongo.Client, catalog string, owner string, name string, scoreType string, timeFrame int) (Score, int) {
	collection := mongoClient.Database("OSS-Score").Collection(catalog) // TODO MAKE DB NAME ENV VAR
	res := GetRepoFromDB(collection, owner, name)
	var score Score
	shelfLife := 7 // Days TODO: make env var
	var repoInfo RepoInfo

	if res.Err() != mongo.ErrNoDocuments { // No match in DB
		err := res.Decode(&repoInfo)

		if err != nil {
			log.Fatalln(err)
		}
		expireDate := time.Now().AddDate(0, 0, -shelfLife)
		if repoInfo.UpdatedAt.After(expireDate) && repoInfo.ScoreStatus == 2 {
			repoWeight := 0.75
			dependencyWeight := 1 - repoWeight
			if scoreType == "activity" {
				score = Score{
					Score:      (repoInfo.RepoActivityScore.Score * repoWeight) + (repoInfo.DependencyActivityScore.Score * dependencyWeight),
					Confidence: (repoInfo.RepoActivityScore.Confidence * repoWeight) + (repoInfo.DependencyActivityScore.Confidence * dependencyWeight),
				}
			} else if scoreType == "license" {
				score = Score{
					Score:      (repoInfo.RepoLicenseScore.Score * repoWeight) + (repoInfo.DependencyLicenseScore.Score * dependencyWeight),
					Confidence: (repoInfo.RepoLicenseScore.Confidence * repoWeight) + (repoInfo.DependencyLicenseScore.Confidence * dependencyWeight),
				}

			}
		}
	}

	return score, repoInfo.ScoreStatus
}

func CalculateScore(catalog string, owner string, name string, timeFrame int, level int) (Score, bool) {
	uri := os.Getenv("MONGO_URI")
	// Create a new mongo_client and connect to the server
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	defer func() {
		if err = mongoClient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	// Ping the primary
	if err := mongoClient.Ping(context.Background(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")

	return calculateScoreHelper(mongoClient, catalog, owner, name, timeFrame, level)
}

func calculateScoreHelper(mongoClient *mongo.Client, catalog string, owner string, name string, timeFrame int, level int) (Score, bool) {
	shelfLife := 7 // Days TODO: make env var
	infoReady := false
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"owner", owner}},
				bson.D{{"name", name}},
			}},
	}
	var repoInfo RepoInfo

	collection := mongoClient.Database("OSS-Score").Collection(catalog) // TODO MAKE DB NAME ENV VAR
	res := GetRepoFromDB(collection, owner, name)

	if res.Err() == mongo.ErrNoDocuments { // No match in DB
		fmt.Println("need to do full query")
		repoInfo = queryGithub(catalog, owner, name, time.Now().AddDate(-(timeFrame/12), -(timeFrame%12), 0))
		infoReady = true // temp while this is synchronous
		repoInfo.ScoreStatus = 1

		fmt.Println("Inserting Data")
		_, err := collection.InsertOne(context.TODO(), repoInfo)
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

			repoInfo := queryGithub(catalog, owner, name, repoInfo.UpdatedAt) // pull only needed data
			infoReady = true                                                  // temp while this is synchronous

			insertableData := bson.D{primitive.E{Key: "$set", Value: repoInfo}}
			repoInfo.ScoreStatus = 1
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
		var wg sync.WaitGroup
		counter := 0
		for _, dependency := range repoInfo.Dependencies {
			if counter < 10 {
				counter += 1
				wg.Add(1)
				go func(catalog string, owner string, name string, timeFrame int, level int) {
					defer wg.Done()
					fmt.Println("updating: " + owner + "/" + name)
					calculateScoreHelper(mongoClient, catalog, owner, name, timeFrame, level)
				}(dependency.Catalog, dependency.Owner, dependency.Name, timeFrame, level)
			} else {
				counter = 0
				wg.Wait()
			}
		}
		wg.Wait()
	}

	repoScore, dependencyScore := CalculateActivityScore(mongoClient, &repoInfo, time.Now().AddDate(-(timeFrame/12), -(timeFrame%12), 0)) // startpoint hardcoded for now

	repoInfo.RepoActivityScore = repoScore
	repoInfo.DependencyActivityScore = dependencyScore
	repoInfo.ScoreStatus = 2

	insertableData := bson.D{primitive.E{Key: "$set", Value: repoInfo}}
	_, err := collection.UpdateOne(context.TODO(), filter, insertableData)
	if err != nil {
		log.Fatal(err)
	}

	repoWeight := 0.75
	dependencyWeight := 1 - repoWeight

	return Score{
		Score:      (repoScore.Score * repoWeight) + (dependencyScore.Score * dependencyWeight),
		Confidence: (repoScore.Confidence * repoWeight) + (dependencyScore.Confidence * dependencyWeight),
	}, infoReady
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

	GetCoreRepoInfo(httpClient, &repoInfo)

	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		GetGithubIssues(httpClient, &repoInfo, startPoint.Format(time.RFC3339))
	}()
	go func() {
		defer wg.Done()
		GetGithubDependencies(httpClient, &repoInfo)
	}()
	go func() {
		defer wg.Done()
		GetGithubCommits(httpClient, &repoInfo, startPoint.Format(time.RFC3339))
	}()
	go func() {
		defer wg.Done()
		GetGithubReleases(httpClient, &repoInfo, startPoint.Format(time.RFC3339))
	}()

	wg.Wait()

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
