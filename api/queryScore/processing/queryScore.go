package main

import (
	"context"
	"fmt"

	runtime "github.com/aws/aws-lambda-go/lambda"
	// "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type message struct {
	Records []types.Message
}

func handler(ctx context.Context, messageOutput types.ReceiveMessageOutput) error {
	fmt.Println(messageOutput)
	// fmt.Println(messageOutput.Messages[0])
	// util.QueryProject(repo.Catalog, repo.Owner, repo.Name, repo.TimeFrame)
	return nil
}

func main() {
	runtime.Start(handler)
}
