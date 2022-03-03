package main

import (
	"api/util"
	"context"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	dbClient := util.GetDynamoDBClient(ctx)
	for _, message := range sqsEvent.Records {
		catalog := *message.MessageAttributes["catalog"].StringValue
		owner := *message.MessageAttributes["owner"].StringValue
		name := *message.MessageAttributes["name"].StringValue
		timeFrame, err := strconv.Atoi(*message.MessageAttributes["timeFrame"].StringValue)

		if err != nil {
			log.Fatalln("Error converting time frame to int")
		}

		util.QueryProject(ctx, dbClient, catalog, owner, name, timeFrame)
	}

	return nil
}

func main() {
	runtime.Start(handler)
}
