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

		repo, err := util.QueryProject(ctx, dbClient, catalog, owner, name, timeFrame)
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println(repo)
	}

	// for _, dependency := range repo.Dependencies {
	// 	fmt.Println("submitting dep to queue")
	// 	util.SubmitDependencies(ctx, dependency.Catalog, dependency.Owner, dependency.Name)
	// }

	return nil
}

func main() {
	runtime.Start(handler)
}
