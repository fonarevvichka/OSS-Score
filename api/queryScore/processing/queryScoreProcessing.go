package main

import (
	"api/util"
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
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

		mongoClient := util.GetMongoClient()
		defer mongoClient.Disconnect(context.TODO())
		collection := mongoClient.Database("OSS-Score").Collection(catalog) // TODO MAKE DB NAME ENV VAR
		util.UpdateScoreState(collection, catalog, owner, name, 2)
		repo, err = util.QueryProject(collection, catalog, owner, name, timeFrame)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	for _, dependency := range repo.Dependencies {
		fmt.Println("submitting dep to queue")
		util.SubmitDependencies(ctx, dependency.Catalog, dependency.Owner, dependency.Name)
	}

	return nil
}

func main() {
	runtime.Start(handler)
}
