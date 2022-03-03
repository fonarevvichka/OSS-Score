package main

import (
	"api/util"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	var repo util.RepoInfo
	for _, message := range sqsEvent.Records {
		catalog := *message.MessageAttributes["catalog"].StringValue
		owner := *message.MessageAttributes["owner"].StringValue
		name := *message.MessageAttributes["name"].StringValue
		timeFrame, err := strconv.Atoi(*message.MessageAttributes["timeFrame"].StringValue)

		if err != nil {
			log.Println("Error converting time frame to int")
			return err
		}

		dbClient := util.GetDynamoDBClient(ctx)
		err = util.SetScoreState(ctx, dbClient, catalog, owner, name, 2)
		if err != nil {
			return err
		}

		repo, err = util.QueryProject(ctx, dbClient, catalog, owner, name, timeFrame)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	queueName := os.Getenv("QUERY_QUEUE")
	sqsClient := util.GetSqsClient(ctx)

	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &queueName,
	}

	result, err := sqsClient.GetQueueUrl(ctx, gQInput)
	if err != nil {
		log.Println("Got an error getting the queue URL:")
		log.Println(err)
		return fmt.Errorf("GetQueueUrl %v", err)
	}

	for _, dependency := range repo.Dependencies {
		fmt.Println("submitting dep to queue")
		util.SubmitDependencies(ctx, sqsClient, *result.QueueUrl, dependency.Catalog, dependency.Owner, dependency.Name)
	}

	return nil
}

func main() {
	runtime.Start(handler)
}
