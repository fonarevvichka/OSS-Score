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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/oauth2"
)

func GetSqsSession(ctx context.Context) *sqs.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	return sqs.NewFromConfig(cfg)
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

func getRepoFilter(owner string, name string) bson.D {
	return bson.D{
		{"$and",
			bson.A{
				bson.D{{"owner", owner}},
				bson.D{{"name", name}},
			}},
	}
}

func getManyRepoFilter(repos []NameOwner) bson.D {
	var filters bson.A
	for _, repo := range repos {
		currFilter := bson.D{
			{"$and",
				bson.A{
					bson.D{{"owner", repo.Owner}},
					bson.D{{"name", repo.Name}},
				}},
		}

		filters = append(filters, currFilter)
	}

	return bson.D{
		{"$or", filters},
	}
}

func GetRepoFromDB(collection *mongo.Collection, owner string, name string) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return collection.FindOne(ctx, getRepoFilter(owner, name))
}

func GetReposFromDB(collection *mongo.Collection, repos []NameOwner) []RepoInfo {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, getManyRepoFilter(repos))

	if err != nil {
		log.Fatalln(err)
	}

	var deps []RepoInfo

	for cur.Next(context.TODO()) {
		var dep RepoInfo
		err := cur.Decode(&dep)

		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		deps = append(deps, dep)
	}

	return deps
}

func GetScore(collection *mongo.Collection, catalog string, owner string, name string, scoreType string, timeFrame int) (Score, int) {
	res := GetRepoFromDB(collection, owner, name)

	var combinedScore Score
	var repoScore Score
	var depScore Score

	shelfLife, err := strconv.Atoi(os.Getenv("SHELF_LIFE"))
	if err != nil {
		log.Fatalln(err)
	}
	var repoInfo RepoInfo

	if res.Err() != mongo.ErrNoDocuments { // Match in DB
		err := res.Decode(&repoInfo)

		if err != nil {
			log.Fatalln(err)
		}
		expireDate := time.Now().AddDate(0, 0, -shelfLife)
		startPoint := time.Now().AddDate(-(timeFrame / 12), -(timeFrame % 12), 0)

		if repoInfo.UpdatedAt.After(expireDate) && repoInfo.Status == 3 {
			repoWeight := 0.75
			dependencyWeight := 1 - repoWeight
			if scoreType == "activity" {
				repoScore = CalculateActivityScore(&repoInfo, startPoint)
				depScore = CalculateDependencyActivityScore(collection, &repoInfo, startPoint)
			} else if scoreType == "license" {
				licenseMap := GetLicenseMap()
				repoScore = CalculateLicenseScore(&repoInfo, licenseMap)
				depScore = CalculateDependencyLicenseScore(collection, &repoInfo, licenseMap)
			}
			combinedScore = Score{
				Score:      (repoScore.Score * repoWeight) + (depScore.Score * dependencyWeight),
				Confidence: (repoScore.Confidence * repoWeight) + (depScore.Confidence * dependencyWeight),
			}
		}
	}

	return combinedScore, repoInfo.Status
}

// Channel returns: RepoInfo struct, dataStatus int -- (0 - nothing new, 1 - updated, 2 - all new data)
func addUpdateRepo(collection *mongo.Collection, catalog string, owner string, name string, timeFrame int, licenseMap map[string]int) RepoInfoMessage {
	shelfLife, err := strconv.Atoi(os.Getenv("SHELF_LIFE"))
	if err != nil {
		log.Fatalln(err)
	}

	res := GetRepoFromDB(collection, owner, name)

	repoInfo := RepoInfo{}
	insert := false

	startPoint := time.Now().AddDate(-(timeFrame / 12), -(timeFrame % 12), 0)

	if res.Err() == mongo.ErrNoDocuments { // No data on repo
		insert = true
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

			log.Println(owner + "/" + name + " Need to do partial update")
			if repoInfo.UpdatedAt.After(startPoint) {
				startPoint = repoInfo.UpdatedAt
			}
			repoInfo = QueryGithub(catalog, owner, name, startPoint) // pull only needed data
			log.Println(owner + "/" + name + " Done querying github")
		}
	}

	return RepoInfoMessage{
		RepoInfo: repoInfo,
		Insert:   insert,
	}
}

func GetLicenseMap() map[string]int {
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

	return licenseMap
}

func SubmitDependencies(ctx context.Context, catalog string, owner string, name string) error {
	queueName := os.Getenv("QUERY_QUEUE")
	fmt.Println(queueName)
	client := GetSqsSession(ctx)

	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &queueName,
	}

	result, err := client.GetQueueUrl(ctx, gQInput)
	if err != nil {
		log.Println("Got an error getting the queue URL:")
		log.Println(err)
		return err
	}

	queueURL := result.QueueUrl
	timeFrame := "12"

	sMInput := &sqs.SendMessageInput{
		MessageAttributes: map[string]types.MessageAttributeValue{
			"catalog": {
				DataType:    aws.String("String"),
				StringValue: &catalog,
			},
			"owner": {
				DataType:    aws.String("String"),
				StringValue: &owner,
			},
			"name": {
				DataType:    aws.String("String"),
				StringValue: &name,
			},
			"timeFrame": {
				DataType:    aws.String("String"),
				StringValue: &timeFrame, // temp hardcoded
			},
		},
		MessageBody: aws.String("Repo to be queried"),
		QueueUrl:    queueURL,
	}

	_, err = client.SendMessage(ctx, sMInput)
	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		return err
	}

	return nil
}

func UpdateScoreState(collection *mongo.Collection, catalog string, owner string, name string, status int) {
	res := GetRepoFromDB(collection, owner, name)

	var repo RepoInfo
	new := false

	if res.Err() == mongo.ErrNoDocuments { // No match in DB
		repo = RepoInfo{ // this is bad, will break, need to get repo from db first
			Catalog: catalog,
			Owner:   owner,
			Name:    name,
			Status:  status,
		}
		new = true
	} else {
		err := res.Decode(&repo)
		if err != nil {
			log.Fatalln(err)
		}
		repo.Status = status
	}

	syncRepoWithDB(collection, repo, new)
}

func QueryProject(collection *mongo.Collection, catalog string, owner string, name string, timeFrame int) RepoInfo {
	licenseMap := GetLicenseMap()

	// get repo info message
	repoInfoMessage := addUpdateRepo(collection, catalog, owner, name, timeFrame, licenseMap)

	mainRepo := repoInfoMessage.RepoInfo
	insert := repoInfoMessage.Insert

	mainRepo.Status = 3
	syncRepoWithDB(collection, mainRepo, insert)

	return mainRepo
}

func syncRepoWithDB(collection *mongo.Collection, repo RepoInfo, new bool) {
	if new {
		_, err := collection.InsertOne(context.TODO(), repo)
		if err != nil {
			log.Println(repo.Owner + "/" + repo.Name)
			log.Fatal(err)
		}
	} else {
		insertableData := bson.D{primitive.E{Key: "$set", Value: repo}}
		filter := getRepoFilter(repo.Owner, repo.Name)
		_, err := collection.UpdateOne(context.TODO(), filter, insertableData)
		if err != nil {
			log.Fatal(err)
		}
	}
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
