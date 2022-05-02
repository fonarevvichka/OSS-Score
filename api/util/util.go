package util

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
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
	"golang.org/x/sync/errgroup"
)

func GetSqsClient(ctx context.Context) *sqs.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	return sqs.NewFromConfig(cfg)
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

func getRepoFilter(owner string, name string) bson.D {
	return bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.M{"owner": owner},
				bson.M{"name": name},
			}},
	}
}

func getManyRepoFilter(repos []NameOwner) bson.M {
	var filters bson.A
	for _, repo := range repos {
		currFilter := bson.D{
			{Key: "$and",
				Value: bson.A{
					bson.M{"owner": repo.Owner},
					bson.M{"name": repo.Name},
				}},
		}

		filters = append(filters, currFilter)
	}

	return bson.M{"$or": filters}
}

func GetRepoFromDB(ctx context.Context, collection *mongo.Collection, catalog string, owner string, name string) (RepoInfo, bool, error) {
	var repo RepoInfo
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	res := collection.FindOne(ctx, getRepoFilter(owner, name))

	if res.Err() != mongo.ErrNoDocuments { // FOUND REPO
		err := res.Decode(&repo)

		if err != nil {
			log.Println(err)
			return repo, false, fmt.Errorf("SingleResult.Decode: %v", err)
		}
	} else { // NOT FOUND
		return RepoInfo{
			Catalog: catalog,
			Owner:   owner,
			Name:    name,
		}, false, nil
	}

	return repo, true, nil
}

func GetReposFromDB(ctx context.Context, collection *mongo.Collection, repos []NameOwner) ([]RepoInfo, error) {
	var deps []RepoInfo

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, getManyRepoFilter(repos))

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

func GetScore(ctx context.Context, collection *mongo.Collection, catalog string, owner string, name string, scoreType string, timeFrame int) (Score, float64, int, string, error) {
	var combinedScore Score
	var repoScore Score
	var depScore Score
	var depRatio float64
	var message string
	repoWeight := 0.75
	dependencyWeight := 1 - repoWeight

	repo, found, err := GetRepoFromDB(ctx, collection, catalog, owner, name)
	if err != nil {
		return combinedScore, depRatio, repo.Status, "", err
	}

	shelfLife, err := strconv.Atoi(os.Getenv("SHELF_LIFE"))
	if err != nil {
		return combinedScore, depRatio, repo.Status, "", err
	}

	if found { // Match in DB
		expireDate := time.Now().AddDate(0, 0, -shelfLife)
		startPoint := time.Now().AddDate(-(timeFrame / 12), -(timeFrame % 12), 0)

		if startPoint.Before(repo.CreatedAt) {
			startPoint = repo.CreatedAt
		}

		if repo.Status == 3 {
			if repo.UpdatedAt.After(expireDate) && (repo.DataStartPoint.Before(startPoint) || repo.DataStartPoint.Equal(startPoint)) {
				if scoreType == "activity" {
					repoScore, err = CalculateRepoActivityScore(&repo, startPoint)
					if err != nil {
						log.Println(err)
						return combinedScore, depRatio, repo.Status, "", err
					}
					// depScore, depRatio, err = CalculateDependencyActivityScore(ctx, collection, &repo, startPoint)
					if err != nil {
						log.Println(err)
						return combinedScore, depRatio, repo.Status, "", err
					}
				} else if scoreType == "license" {
					licenseMap, err := GetLicenseMap("./util/scores/licenseScoring.csv")
					if err != nil {
						log.Println(err)
						return combinedScore, depRatio, repo.Status, "", err
					}
					repoScore = CalculateRepoLicenseScore(&repo, licenseMap)
					// depScore, depRatio, err = CalculateDependencyLicenseScore(ctx, collection, &repo, licenseMap)
					if err != nil {
						return combinedScore, depRatio, repo.Status, "", err
					}
				}

				// if there are no deps we want to not include them in the score
				// TEMP MEASURE TO NOT INCLUDE THE SCORES OF DEPENDENCIES
				if len(repo.Dependencies) == 0 || true {
					repoWeight = 1
					dependencyWeight = 0
				}

				combinedScore = Score{
					Score:      (repoScore.Score * repoWeight) + (depScore.Score * dependencyWeight),
					Confidence: (repoScore.Confidence * repoWeight) + (depScore.Confidence * dependencyWeight),
				}
			} else {
				message = "Data out of date"
			}
		}
	}

	// HARDCODE DEP RATIO TO 1 to prevent chrome extension from spamming
	depRatio = 1
	return combinedScore, depRatio, repo.Status, message, nil
}

func addUpdateRepo(ctx context.Context, collection *mongo.Collection, catalog string, owner string, name string, timeFrame int) (RepoInfo, error) {
	shelfLife, err := strconv.Atoi(os.Getenv("SHELF_LIFE"))
	if err != nil {
		log.Println(err)
		return RepoInfo{}, err
	}

	repo, found, err := GetRepoFromDB(ctx, collection, catalog, owner, name)

	if err != nil {
		return RepoInfo{}, err
	}

	startPoint := time.Now().AddDate(-(timeFrame / 12), -(timeFrame % 12), 0)

	if !found { // No data on repo, essentially useless, since it will be inserted by handler
		log.Println(owner + "/" + name + " Not in DB, need to do full query")
		repo.DataStartPoint = startPoint

		err = QueryGithub(ctx, &repo, startPoint)

		if err != nil {
			log.Println(err)
			return repo, err
		}
		log.Println(owner + "/" + name + " Done querying github")
	} else {
		// Repo data expired, or empty
		if repo.UpdatedAt.Before(time.Now().AddDate(0, 0, -shelfLife)) || startPoint.Before(repo.DataStartPoint) {
			log.Println(owner + "/" + name + " Need to query github")

			if startPoint.Before(repo.UpdatedAt) && repo.DataStartPoint.Before(startPoint) { // set start point to collect only needed data
				startPoint = repo.UpdatedAt
			}

			err = QueryGithub(ctx, &repo, startPoint)
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

	if repo.DataStartPoint.Before(repo.CreatedAt) {
		repo.DataStartPoint = repo.CreatedAt
	}

	repo.UpdatedAt = time.Now()
	return repo, nil
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

func QueryProject(ctx context.Context, collection *mongo.Collection, catalog string, owner string, name string, timeFrame int) (RepoInfo, error) {
	repo, err := addUpdateRepo(ctx, collection, catalog, owner, name, timeFrame)

	if err != nil {
		log.Println(err)
		return RepoInfo{}, err
	}

	repo.Status = 3
	err = SyncRepoWithDB(ctx, collection, repo)

	return repo, err
}

func SetScoreState(ctx context.Context, collection *mongo.Collection, catalog string, owner string, name string, status int) error {
	// setting catalog here to be safe in case of upsert
	insertableData := bson.D{primitive.E{Key: "$set", Value: bson.M{"status": status, "catalog": catalog}}}
	filter := getRepoFilter(owner, name)
	upsert := true

	opts := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := collection.UpdateOne(ctx, filter, insertableData, &opts)
	if err != nil {
		log.Printf("collection.UpdateOne: %v", err)
		return fmt.Errorf("collection.UpdateOne: %v", err)
	}

	return nil
}

func SyncRepoWithDB(ctx context.Context, collection *mongo.Collection, repo RepoInfo) error {
	insertableData := bson.D{primitive.E{Key: "$set", Value: repo}}
	filter := getRepoFilter(repo.Owner, repo.Name)
	upsert := true

	opts := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := collection.UpdateOne(ctx, filter, insertableData, &opts)
	if err != nil {
		log.Printf("collection.UpdateOne: %v", err)
		return fmt.Errorf("collection.UpdateOne: %v", err)
	}

	return nil
}

func QueryGithub(ctx context.Context, repo *RepoInfo, startPoint time.Time) error {
	errs, ctx := errgroup.WithContext(ctx)

	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_PAT")},
	))

	errs.Go(func() error {
		return GetGithubIssuesRest(ctx, httpClient, repo, startPoint)
	})

	errs.Go(func() error {
		// TODO, not taking into account the timeframe
		return GetGithubReleasesGraphQLManual(httpClient, repo, startPoint.Format(time.RFC3339))
	})

	errs.Go(func() error {
		return GetCoreRepoInfo(httpClient, repo)
	})

	errs.Go(func() error {
		// COULD BE GRAPHQL with 70% performace
		return GetGithubCommitsRest(ctx, httpClient, repo, startPoint)
	})

	errs.Go(func() error {
		return GetGithubPullsGraphQL(httpClient, repo, startPoint)
	})

	errs.Go(func() error {
		return GetGithubDependencies(httpClient, repo)
	})

	return errs.Wait()
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

func readCsv(path string) (data [][]string, err error) {
	file, err := os.Open(path)

	if err != nil {
		return data, fmt.Errorf("os.Open: %v", err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	data, err = csvReader.ReadAll()
	if err != nil {
		return data, fmt.Errorf("csv.ReadAll: %v", err)
	}

	return data, err
}
