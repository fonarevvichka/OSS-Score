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

func getMongoClient() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	// Create a new mongo_client and connect to the server
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalln(err)
	}

	if err := mongoClient.Ping(context.Background(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")

	return mongoClient
}

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

func GetCachedScore(mongoClient *mongo.Client, catalog string, owner string, name string, scoreType string, timeFrame int) (Score, int) {
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

func AddUpdateRepo(catalog string, owner string, name string, timeFrame int, level int) {
	fmt.Println("KISS: Keep it simple stupid")

	mongoClient := getMongoClient()
	defer mongoClient.Disconnect(context.TODO())

	shelfLife := 7
	repoInfo := RepoInfo{}

	collection := mongoClient.Database("OSS-Score").Collection(catalog) // TODO MAKE DB NAME ENV VAR
	res := GetRepoFromDB(collection, owner, name)

	if res.Err() != mongo.ErrNoDocuments { // No data on repo
		fmt.Println("Not in DB, need to query full history")
	} else {
		err := res.Decode(&repoInfo)

		if err != nil {
			log.Fatalln(err)
		}

		// Repo data out of date, need to update
		if repoInfo.UpdatedAt.Before(time.Now().AddDate(0, 0, -shelfLife)) {
			fmt.Println("Need to do partial update")
		} else {
			fmt.Println("Data up to date")
		}
	}
}

func CalculateScore(catalog string, owner string, name string, timeFrame int, level int) Score {

	return calculateScoreHelper(mongoClient, catalog, owner, name, timeFrame, level)
}

//check if repo is in db if yes --> done
//if
func calculateScoreHelper(mongoClient *mongo.Client, catalog string, owner string, name string, timeFrame int, level int) Score {
	shelfLife := 7 // Days TODO: make env var
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
		fmt.Println(owner + "/" + name + " need to do full query")
		repoInfo = queryGithub(catalog, owner, name, time.Now().AddDate(-(timeFrame/12), -(timeFrame%12), 0))
		fmt.Println(owner + "/" + name + "Done querying github")
		repoInfo.ScoreStatus = 1

		fmt.Println(owner + "/" + name + "Inserting Data")
		_, err := collection.InsertOne(context.TODO(), repoInfo)
		if err != nil {
			fmt.Println(repoInfo.Owner + "/" + repoInfo.Name)
			log.Fatal(err)
		}
		fmt.Println(owner + "/" + name + "Data Inserted")

	} else { // Match in DB found
		err := res.Decode(&repoInfo)

		if err != nil {
			log.Fatalln(err)
		}

		expireDate := time.Now().AddDate(0, 0, -shelfLife)
		if repoInfo.UpdatedAt.Before(expireDate) {
			fmt.Println(repoInfo.Owner + "/" + repoInfo.Name + " out of date: need to make partial query")

			repoInfo := queryGithub(catalog, owner, name, repoInfo.UpdatedAt) // pull only needed data

			insertableData := bson.D{primitive.E{Key: "$set", Value: repoInfo}}
			repoInfo.ScoreStatus = 1
			_, err := collection.UpdateOne(context.TODO(), filter, insertableData)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if level < 1 {
		level += 1
		var wg sync.WaitGroup
		counter := 0
		fmt.Println("Querying dependencies")
		for _, dependency := range repoInfo.Dependencies {
			if counter < 10 { // caps concurrent works at 10
				counter += 1
				wg.Add(1)
				go func(catalog string, owner string, name string, timeFrame int, level int) {
					defer wg.Done()
					calculateScoreHelper(mongoClient, catalog, owner, name, timeFrame, level)
				}(dependency.Catalog, dependency.Owner, dependency.Name, timeFrame, level)
			} else {
				counter = 0
				wg.Wait()
			}
		}
		wg.Wait()
	}

	repoScore := CalculateRepoActivityScore(&repoInfo, time.Now().AddDate(-(timeFrame/12), -(timeFrame%12), 0))
	dependencyScore := CalculateDependencyActivityScore(collection, &repoInfo, time.Now().AddDate(-(timeFrame/12), -(timeFrame%12), 0)) // startpoint hardcoded for now

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
