package main

import (
	"api/util"
	"context"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	// var repo util.RepoInfo
	for _, message := range sqsEvent.Records {
		catalog := *message.MessageAttributes["catalog"].StringValue
		owner := *message.MessageAttributes["owner"].StringValue
		name := *message.MessageAttributes["name"].StringValue
		timeFrame, err := strconv.Atoi(*message.MessageAttributes["timeFrame"].StringValue)

		if err != nil {
			log.Println("Error converting time frame to int")
			return err
		}

		mongoClient, connected, err := util.GetMongoClient(ctx)
		if connected {
			defer mongoClient.Disconnect(ctx)
		}

		if err != nil {
			return err
		}
		collection := mongoClient.Database(os.Getenv("MONGO_DB")).Collection(catalog)

		err = util.SetScoreState(ctx, collection, catalog, owner, name, 2)
		if err != nil {
			return err
		}

		_, err = util.QueryProject(ctx, collection, catalog, owner, name, timeFrame)
		if err != nil {
			log.Println(err)
			util.SetScoreState(ctx, collection, catalog, owner, name, 4)
			return err
		}
	}

	// sqsClient := util.GetSqsClient(ctx)
	// queueURL := os.Getenv("QUEUE_URL")
	// for _, dependency := range repo.Dependencies {
	// 	log.Println("submitting dep to queue")
	// 	util.SubmitDependencies(ctx, sqsClient, queueURL, dependency.Catalog, dependency.Owner, dependency.Name)
	// }

	return nil
}

func main() {
	runtime.Start(handler)
}
