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
	mongoClient := util.GetMongoClient()
	defer mongoClient.Disconnect(context.TODO())
	for _, message := range sqsEvent.Records {
		catalog := *message.MessageAttributes["catalog"].StringValue
		owner := *message.MessageAttributes["owner"].StringValue
		name := *message.MessageAttributes["name"].StringValue
		timeFrame, err := strconv.Atoi(*message.MessageAttributes["timeFrame"].StringValue)

		if err != nil {
			log.Fatalln("Error converting time frame to int")
		}

		collection := mongoClient.Database("OSS-Score").Collection(catalog) // TODO MAKE DB NAME ENV VAR
		util.QueryProject(collection, catalog, owner, name, timeFrame)
	}

	return nil
}

func main() {
	runtime.Start(handler)
}
