package util

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/oauth2"
)

func getRepoFilter(owner string, name string) bson.D {
	return bson.D{
		{"$and",
			bson.A{
				bson.D{{"owner", owner}},
				bson.D{{"name", name}},
			}},
	}
}

func GetMongoClient() *mongo.Client {
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

	return collection.FindOne(ctx, getRepoFilter(owner, name))
}

func GetCachedScore(mongoClient *mongo.Client, catalog string, owner string, name string, scoreType string, timeFrame int) (Score, int) {
	collection := mongoClient.Database("OSS-Score").Collection(catalog) // TODO MAKE DB NAME ENV VAR
	res := GetRepoFromDB(collection, owner, name)
	var score Score
	shelfLife, err := strconv.Atoi(os.Getenv("SHELF_LIFE"))
	if err != nil {
		log.Fatalln(err)
	}
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

// Channel returns: RepoInfo struct, dataStatus int -- (0 - nothing new, 1 - updated, 2 - all new data)
func addUpdateRepo(collection *mongo.Collection, catalog string, owner string, name string, timeFrame int, licenseMap map[string]int) RepoInfoMessage {
	shelfLife, err := strconv.Atoi(os.Getenv("SHELF_LIFE"))
	if err != nil {
		log.Fatalln(err)
	}

	res := GetRepoFromDB(collection, owner, name)

	repoInfo := RepoInfo{}
	dataStatus := 0
	startPoint := time.Now().AddDate(-(timeFrame / 12), -(timeFrame % 12), 0)

	if res.Err() == mongo.ErrNoDocuments { // No data on repo
		dataStatus = 2
		log.Println(owner + "/" + name + " Not in DB, need to do full query")

		repoInfo = QueryGithub(catalog, owner, name, startPoint)
		log.Println(owner + "/" + name + " Done querying github")
	} else {
		err := res.Decode(&repoInfo)
		if err != nil {
			log.Fatalln(err)
		}

		// Repo data expired
		if repoInfo.UpdatedAt.Before(time.Now().AddDate(0, 0, -shelfLife)) {
			dataStatus = 1

			log.Println(owner + "/" + name + " Need to do partial update")
			repoInfo = QueryGithub(catalog, owner, name, repoInfo.UpdatedAt) // pull only needed data
			log.Println(owner + "/" + name + " Done querying github")
		}
	}

	if dataStatus != 0 {
		repoInfo.RepoActivityScore = CalculateRepoActivityScore(&repoInfo, startPoint)
		repoInfo.RepoLicenseScore = CalculateRepoLicenseScore(&repoInfo, licenseMap)
	}

	return RepoInfoMessage{
		RepoInfo:   repoInfo,
		DataStatus: dataStatus,
	}
}

func QueryProject(catalog string, owner string, name string, timeFrame int) {
	mongoClient := GetMongoClient()
	defer mongoClient.Disconnect(context.TODO())
	collection := mongoClient.Database("OSS-Score").Collection(catalog) // TODO MAKE DB NAME ENV VAR

	// Get License Score map
	licenseMap := make(map[string]int)

	licenseFile, err := os.Open("./scores/licenseScores.txt")

	if err != nil {
		log.Fatalln(err)
	}

	defer licenseFile.Close()

	scanner := bufio.NewScanner(licenseFile)

	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Split(line, ",")
		score, err := strconv.Atoi(values[1])
		if err != nil {
			log.Fatalln(err)
		}
		licenseMap[values[0]] = score
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// get repo info message
	repoInfoMessage := addUpdateRepo(collection, catalog, owner, name, timeFrame, licenseMap)

	mainRepo := repoInfoMessage.RepoInfo
	dataStatus := repoInfoMessage.DataStatus

	updateDependencies(collection, &mainRepo, timeFrame, licenseMap)

	startPoint := time.Now().AddDate(-(timeFrame / 12), -(timeFrame % 12), 0)
	mainRepo.DependencyActivityScore = CalculateDependencyActivityScore(collection, &mainRepo, startPoint)
	mainRepo.ScoreStatus = 2
	mainRepo.DependencyLicenseScore = CalculateDependencyLicenseScore(collection, &mainRepo)
	syncRepoWithDB(collection, mainRepo, dataStatus)

	log.Println("DONE!")
}

func syncRepoWithDB(collection *mongo.Collection, repo RepoInfo, dataStatus int) {
	if dataStatus == 1 {
		insertableData := bson.D{primitive.E{Key: "$set", Value: repo}}
		// log.Println(owner + "/" + name + " Updating Data")
		filter := getRepoFilter(repo.Owner, repo.Name)
		_, err := collection.UpdateOne(context.TODO(), filter, insertableData)
		if err != nil {
			log.Fatal(err)
		}
	} else if dataStatus == 2 {
		_, err := collection.InsertOne(context.TODO(), repo)
		if err != nil {
			log.Println(repo.Owner + "/" + repo.Name)
			log.Fatal(err)
		}
	}
}

func updateDependencies(collection *mongo.Collection, mainRepo *RepoInfo, timeFrame int, licenseMap map[string]int) {
	var wg sync.WaitGroup
	var repoMessages []RepoInfoMessage
	dependencies := mainRepo.Dependencies

	counter := 0
	for _, dependency := range dependencies {
		if counter == 1 {
			counter = 0
			wg.Wait()
		} else {
			if dependency.Owner != "" && dependency.Name != "" {
				wg.Add(1)

				go func(collection *mongo.Collection, catalog string, owner string, name string, timeFrame int) {
					defer wg.Done()
					repoMessages = append(repoMessages, addUpdateRepo(collection, catalog, owner, name, timeFrame, licenseMap))
				}(collection, dependency.Catalog, dependency.Owner, dependency.Name, timeFrame)
				counter += 1
			}
		}

		// cap on how many deps to query: testing only
		// if counter == 25 {
		// 	break
		// }
	}
	wg.Wait()

	var newDeps []interface{}
	var updatedDeps []RepoInfo

	for _, repoMessage := range repoMessages {
		if repoMessage.DataStatus == 1 {
			updatedDeps = append(updatedDeps, repoMessage.RepoInfo)
		} else if repoMessage.DataStatus == 2 {
			newDeps = append(newDeps, repoMessage.RepoInfo)
		}
	}

	log.Println("Done querying all deps")
	//TODO: Investigate BulkWrite operation
	if len(newDeps) != 0 {
		log.Println("inserting new deps")
		wg.Add(1)
		go func(collection *mongo.Collection, deps []interface{}) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			collection.InsertMany(ctx, newDeps)
		}(collection, newDeps)
	}

	for _, dep := range updatedDeps {
		wg.Add(1)
		//TODO: Should use syncRepoWithDB
		go func(collection *mongo.Collection, dep RepoInfo) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			collection.UpdateOne(ctx, getRepoFilter(dep.Owner, dep.Name), dep)
		}(collection, dep)

	}

	wg.Wait()
}

func QueryGithub(catalog string, owner string, name string, startPoint time.Time) RepoInfo {
	src1 := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_PAT_1")},
	)
	src2 := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_PAT_2")},
	)
	src3 := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_PAT_3")},
	)
	src4 := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_PAT_4")},
	)
	src5 := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_PAT_5")},
	)

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
	httpClient1 := oauth2.NewClient(context.Background(), src1)
	httpClient2 := oauth2.NewClient(context.Background(), src2)
	httpClient3 := oauth2.NewClient(context.Background(), src3)
	httpClient4 := oauth2.NewClient(context.Background(), src4)
	httpClient5 := oauth2.NewClient(context.Background(), src5)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		GetGithubIssuesRest(httpClient1, &repoInfo, startPoint.Format(time.RFC3339))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		GetGithubDependencies(httpClient2, &repoInfo)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		GetGithubReleases(httpClient3, &repoInfo, startPoint.Format(time.RFC3339))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		GetCoreRepoInfo(httpClient4, &repoInfo)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		GetGithubCommitsRest(httpClient5, &repoInfo, startPoint.Format(time.RFC3339))
	}()

	wg.Wait()

	return repoInfo
}

func dependencyInSlice(dependency Dependency, dependencies []Dependency) bool {
	for _, elem := range dependencies {
		if dependency.Catalog == elem.Catalog &&
			dependency.Owner == elem.Owner &&
			dependency.Name == elem.Name { // NOT SURE IF DEEP COMPARE LIKE THIS IS NEEDED
			return true
		}
	}
	return false
}
