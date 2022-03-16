package util

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

func GetSqsClient(ctx context.Context) *sqs.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	return sqs.NewFromConfig(cfg)
}

func GetDynamoDBClient(ctx context.Context) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(cfg)
}

func GetMongoClient(ctx context.Context) (*mongo.Client, bool, error) {
	uri := os.Getenv("MONGO_URI")
	// Create a new mongo_client and connect to the server
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPIOptions)

	mongoClient, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return mongoClient, false, fmt.Errorf("mongo.Connect: %v", err)
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return mongoClient, true, fmt.Errorf("mongo.Ping: %v", err)
	}
	fmt.Println("Successfully connected and pinged.")

	return mongoClient, true, nil
}

func getRepoFilterMongo(owner string, name string) bson.D {
	return bson.D{
		{"$and",
			bson.A{
				bson.D{{"owner", owner}},
				bson.D{{"name", name}},
			}},
	}
}

func getManyRepoFilterMongo(repos []NameOwner) bson.D {
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

func GetRepoFromDBMongo(ctx context.Context, collection *mongo.Collection, owner string, name string) (RepoInfo, bool, error) {
	var repo RepoInfo
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	res := collection.FindOne(ctx, getRepoFilterMongo(owner, name))

	if res.Err() != mongo.ErrNoDocuments { // FOUND REPO
		err := res.Decode(&repo)

		if err != nil {
			log.Println(err)
			return repo, false, fmt.Errorf("SingleResult.Decode: %v", err)
		}
	} else { // NOT FOUND
		return repo, false, nil
	}

	return repo, true, nil
}

func GetReposFromDBMongo(ctx context.Context, collection *mongo.Collection, repos []NameOwner) ([]RepoInfo, error) {
	var deps []RepoInfo

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, getManyRepoFilterMongo(repos))

	if err != nil {
		log.Println("Error on finding documents", err)
		return deps, fmt.Errorf("collection.Find: %v", err)
	}

	for cur.Next(ctx) {
		var dep RepoInfo
		err := cur.Decode(&dep)

		if err != nil {
			log.Println("Error on Decoding the document", err)
			return deps, fmt.Errorf("SingleResult.Decode: %v", err)
		}
		deps = append(deps, dep)
	}

	return deps, nil
}

// repo, found, error
func GetRepoFromDB(ctx context.Context, client *dynamodb.Client, owner string, name string) (RepoInfo, bool, error) {
	var repo RepoInfo
	data, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
		Key: map[string]dynamoTypes.AttributeValue{
			"name":  &dynamoTypes.AttributeValueMemberS{Value: name},
			"owner": &dynamoTypes.AttributeValueMemberS{Value: owner},
		},
	})

	if err != nil {
		return repo, false, fmt.Errorf("GetItem: %v", err)
	}

	if data.Item == nil {
		return repo, false, nil
	}

	err = attributevalue.UnmarshalMap(data.Item, &repo)
	if err != nil {
		return repo, false, fmt.Errorf("UnmarhsalMap: %v", err)
	}

	return repo, true, nil
}

func GetReposFromDB(ctx context.Context, client *dynamodb.Client, repoKeys []NameOwner) ([]RepoInfo, error) {
	table := os.Getenv("DYNAMODB_TABLE")
	repoKeyChunks := make([][]NameOwner, (len(repoKeys))/100+1)
	repoInfoChunks := make([][]RepoInfo, (len(repoKeys))/100+1)

	chunkCounter := 0
	for i, repoKey := range repoKeys {
		if i%100 == 0 && i/100 > 0 { //Chunk complete
			chunkCounter++
		}
		repoKeyChunks[chunkCounter] = append(repoKeyChunks[chunkCounter], repoKey)
	}

	var errs errgroup.Group

	for i, repoKeyChunk := range repoKeyChunks {
		func(chunk []NameOwner, chunkIndex int) {
			errs.Go(func() error {
				var keys []map[string]dynamoTypes.AttributeValue
				for _, repoKey := range chunk {
					keys = append(keys, map[string]dynamoTypes.AttributeValue{
						"name":  &dynamoTypes.AttributeValueMemberS{Value: repoKey.Name},
						"owner": &dynamoTypes.AttributeValueMemberS{Value: repoKey.Owner},
					})
				}

				data, err := client.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
					RequestItems: map[string]dynamoTypes.KeysAndAttributes{
						table: {
							Keys: keys,
						},
					},
				})
				if err != nil {
					return fmt.Errorf("BatchGetItem: %v", err)
				}
				items := data.Responses[table]

				for _, item := range items {
					var repo RepoInfo
					if item != nil {
						err := attributevalue.UnmarshalMap(item, &repo)
						if err != nil {
							return fmt.Errorf("UnmarhsalMap: %v", err)
						}
						repoInfoChunks[chunkIndex] = append(repoInfoChunks[chunkIndex], repo)
					}
				}
				return nil
			})
		}(repoKeyChunk, i)
	}
	err := errs.Wait()
	var repos []RepoInfo

	for _, repoChunk := range repoInfoChunks {
		repos = append(repos, repoChunk...)
	}

	return repos, err
}

func GetScore(ctx context.Context, collection *mongo.Collection, catalog string, owner string, name string, scoreType string, timeFrame int) (Score, float64, int) {
	repoInfo, found, err := GetRepoFromDBMongo(ctx, collection, owner, name)
	if err != nil {
		log.Fatalln(err)
	}

	var combinedScore Score
	var repoScore Score
	var depScore Score
	var depRatio float64

	shelfLife, err := strconv.Atoi(os.Getenv("SHELF_LIFE"))
	if err != nil {
		log.Fatalln(err)
	}

	if found { // Match in DB
		expireDate := time.Now().AddDate(0, 0, -shelfLife)
		startPoint := time.Now().AddDate(-(timeFrame / 12), -(timeFrame % 12), 0)

		if repoInfo.UpdatedAt.After(expireDate) && repoInfo.Status == 3 {
			repoWeight := 0.75
			dependencyWeight := 1 - repoWeight
			if scoreType == "activity" {
				repoScore = CalculateActivityScore(&repoInfo, startPoint)
				depScore, depRatio, err = CalculateDependencyActivityScoreMongo(ctx, collection, &repoInfo, startPoint)
				log.Println(err)
			} else if scoreType == "license" {
				licenseMap := GetLicenseMap()
				repoScore = CalculateLicenseScore(&repoInfo, licenseMap)
				depScore, depRatio, err = CalculateDependencyLicenseScoreMongo(ctx, collection, &repoInfo, licenseMap)
				log.Println(err)
			}
			combinedScore = Score{
				Score:      (repoScore.Score * repoWeight) + (depScore.Score * dependencyWeight),
				Confidence: (repoScore.Confidence * repoWeight) + (depScore.Confidence * dependencyWeight),
			}
		}
	}

	return combinedScore, depRatio, repoInfo.Status
}

func addUpdateRepo(ctx context.Context, dbClient *dynamodb.Client, catalog string, owner string, name string, timeFrame int) (RepoInfo, error) {
	shelfLife, err := strconv.Atoi(os.Getenv("SHELF_LIFE"))
	if err != nil {
		log.Println(err)
		return RepoInfo{}, err
	}

	repo, found, err := GetRepoFromDB(ctx, dbClient, owner, name)

	if err != nil {
		return RepoInfo{}, err
	}

	startPoint := time.Now().AddDate(-(timeFrame / 12), -(timeFrame % 12), 0)

	if !found { // No data on repo, essentially useless, since it will be inserted by handler
		log.Println(owner + "/" + name + " Not in DB, need to do full query")
		repo.DataStartPoint = startPoint

		repo, err = QueryGithub(catalog, owner, name, startPoint)
		if err != nil {
			log.Println(err)
			return repo, err
		}
		log.Println(owner + "/" + name + " Done querying github")
	} else {
		// Repo data expired, or empty
		if repo.UpdatedAt.Before(time.Now().AddDate(0, 0, -shelfLife)) || startPoint.Before(repo.DataStartPoint) {
			log.Println(owner + "/" + name + " Need to query github")

			if startPoint.Before(repo.UpdatedAt) { // set start point to collect only needed data
				startPoint = repo.UpdatedAt
			}

			repo, err = QueryGithub(catalog, owner, name, startPoint)
			if err != nil {
				log.Println(err)
				return repo, err
			}

			if repo.DataStartPoint.IsZero() || startPoint.Before(repo.DataStartPoint) {
				repo.DataStartPoint = startPoint
			}
			log.Println(owner + "/" + name + " Done querying github")
		}
	}

	return repo, nil
}

func addUpdateRepoMongo(ctx context.Context, collection *mongo.Collection, catalog string, owner string, name string, timeFrame int) (RepoInfo, error) {
	shelfLife, err := strconv.Atoi(os.Getenv("SHELF_LIFE"))
	if err != nil {
		log.Println(err)
		return RepoInfo{}, err
	}

	repo, found, err := GetRepoFromDBMongo(ctx, collection, owner, name)

	if err != nil {
		return RepoInfo{}, err
	}

	startPoint := time.Now().AddDate(-(timeFrame / 12), -(timeFrame % 12), 0)

	if !found { // No data on repo, essentially useless, since it will be inserted by handler
		log.Println(owner + "/" + name + " Not in DB, need to do full query")
		repo.DataStartPoint = startPoint

		repo, err = QueryGithub(catalog, owner, name, startPoint)
		if err != nil {
			log.Println(err)
			return repo, err
		}
		log.Println(owner + "/" + name + " Done querying github")
	} else {
		// Repo data expired, or empty
		if repo.UpdatedAt.Before(time.Now().AddDate(0, 0, -shelfLife)) || startPoint.Before(repo.DataStartPoint) {
			log.Println(owner + "/" + name + " Need to query github")

			if startPoint.Before(repo.UpdatedAt) { // set start point to collect only needed data
				startPoint = repo.UpdatedAt
			}

			repo, err = QueryGithub(catalog, owner, name, startPoint)
			if err != nil {
				log.Println(err)
				return repo, err
			}

			if repo.DataStartPoint.IsZero() || startPoint.Before(repo.DataStartPoint) {
				repo.DataStartPoint = startPoint
			}
			log.Println(owner + "/" + name + " Done querying github")
		}
	}

	return repo, nil
}

func GetLicenseMap() map[string]int {
	// Get License Score map
	licenseMap := make(map[string]int)

	licenseFile, err := os.Open("./util/scores/licenseScores.txt")

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

func SubmitDependencies(ctx context.Context, client *sqs.Client, queueURL string, catalog string, owner string, name string) error {
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
		QueueUrl:    &queueURL,
	}

	_, err := client.SendMessage(ctx, sMInput)
	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		return fmt.Errorf("sqs.client SendMessage %v", err)
	}

	return nil
}

func SetScoreStateMongo(ctx context.Context, collection *mongo.Collection, catalog string, owner string, name string, status int) error {
	repo := RepoInfo{
		Catalog: catalog,
		Owner:   owner,
		Name:    name,
		Status:  status,
	}

	syncRepoWithDBMongo(ctx, collection, repo)

	return nil
}

func SetScoreState(ctx context.Context, dbClient *dynamodb.Client, catalog string, owner string, name string, status int) error {
	repo, found, err := GetRepoFromDB(ctx, dbClient, owner, name)

	if err != nil {
		return err
	}

	if !found { // No match in DB
		repo = RepoInfo{
			Catalog: catalog,
			Owner:   owner,
			Name:    name,
			Status:  status,
		}
	} else {
		repo.Status = status
	}

	_, err = dbClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
		Key: map[string]dynamoTypes.AttributeValue{
			"name":  &dynamoTypes.AttributeValueMemberS{Value: repo.Name},
			"owner": &dynamoTypes.AttributeValueMemberS{Value: repo.Owner},
		},
		UpdateExpression: aws.String("set #catalog = :catalog, #status = :status"),
		ExpressionAttributeValues: map[string]dynamoTypes.AttributeValue{
			":catalog": &dynamoTypes.AttributeValueMemberS{Value: repo.Catalog},
			":status":  &dynamoTypes.AttributeValueMemberN{Value: strconv.Itoa(repo.Status)},
		},
		ExpressionAttributeNames: map[string]string{
			"#catalog": "catalog",
			"#status":  "status",
		},
	})

	if err != nil {
		return fmt.Errorf("Error updating item %v", err)
	}

	return nil
}

func QueryProject(ctx context.Context, dbClient *dynamodb.Client, catalog string, owner string, name string, timeFrame int) (RepoInfo, error) {
	repo, err := addUpdateRepo(ctx, dbClient, catalog, owner, name, timeFrame)

	if err != nil {
		log.Println(err)
		return RepoInfo{}, err
	}

	repo.Status = 3
	err = syncRepoWithDB(ctx, dbClient, repo)

	return repo, err
}

func QueryProjectMongo(ctx context.Context, collection *mongo.Collection, catalog string, owner string, name string, timeFrame int) (RepoInfo, error) {
	repo, err := addUpdateRepoMongo(ctx, collection, catalog, owner, name, timeFrame)

	if err != nil {
		log.Println(err)
		return RepoInfo{}, err
	}

	repo.Status = 3
	err = syncRepoWithDBMongo(ctx, collection, repo)

	return repo, err
}

func syncRepoWithDBMongo(ctx context.Context, collection *mongo.Collection, repo RepoInfo) error {
	insertableData := bson.D{primitive.E{Key: "$set", Value: repo}}
	filter := getRepoFilterMongo(repo.Owner, repo.Name)
	upsert := true

	opts := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := collection.UpdateOne(ctx, filter, insertableData, &opts)
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("collection.UpdateOne: %v", err)
	}

	return nil
}

func syncRepoWithDB(ctx context.Context, client *dynamodb.Client, repo RepoInfo) error {
	data, err := attributevalue.MarshalMap(repo)
	if err != nil {
		return fmt.Errorf("MarshalMap: %v", err)
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE")),
		Item:      data,
	})

	if err != nil {
		log.Println("Error inserting " + repo.Owner + "/" + repo.Name)
		return fmt.Errorf("PutItem: %v", err)
	}

	return nil
}

func QueryGithub(catalog string, owner string, name string, startPoint time.Time) (RepoInfo, error) {
	errs, ctx := errgroup.WithContext(context.Background())

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
	httpClient1 := oauth2.NewClient(ctx, src1)
	httpClient2 := oauth2.NewClient(ctx, src2)
	httpClient3 := oauth2.NewClient(ctx, src3)
	httpClient4 := oauth2.NewClient(ctx, src4)
	httpClient5 := oauth2.NewClient(ctx, src5)

	errs.Go(func() error {
		return GetGithubIssuesRest(httpClient1, &repoInfo, startPoint.Format(time.RFC3339))
	})

	errs.Go(func() error {
		return GetGithubDependencies(httpClient2, &repoInfo)
	})

	errs.Go(func() error {
		return GetGithubReleases(httpClient3, &repoInfo, startPoint.Format(time.RFC3339))
	})

	errs.Go(func() error {
		return GetCoreRepoInfo(httpClient4, &repoInfo)
	})

	errs.Go(func() error {
		return GetGithubCommitsRest(httpClient5, &repoInfo, startPoint.Format(time.RFC3339))
	})

	err := errs.Wait()
	return repoInfo, err
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
