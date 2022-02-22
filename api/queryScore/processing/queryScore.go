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
			log.Fatalln("Error converting time frame to int")
		}

		repo = util.QueryProject(catalog, owner, name, timeFrame, ctx)
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
