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
	mongoClient, connected, err := util.GetMongoClient(ctx)
	if connected {
		defer mongoClient.Disconnect(ctx)
	}
	if err != nil {
		log.Println(err)
		return err
	}

	for _, message := range sqsEvent.Records {
		catalog := *message.MessageAttributes["catalog"].StringValue
		owner := *message.MessageAttributes["owner"].StringValue
		name := *message.MessageAttributes["name"].StringValue
		timeFrame, err := strconv.Atoi(*message.MessageAttributes["timeFrame"].StringValue)

		collection := mongoClient.Database(os.Getenv("MONGO_DB")).Collection(catalog) //FIND OUT IF THIS IS SLOW

		if err != nil {
			log.Fatalln("Error converting time frame to int")
		}

		util.QueryProject(ctx, collection, catalog, owner, name, timeFrame)
		if err != nil {
			log.Println(err)
			util.SetScoreState(ctx, collection, owner, name, 4)
			return err
		}
	}

	return nil
}

func main() {
	runtime.Start(handler)
}
