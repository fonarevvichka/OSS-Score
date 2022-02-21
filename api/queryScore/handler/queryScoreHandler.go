package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type response struct {
	Message string
}

func sqsSession(ctx context.Context) *sqs.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	return sqs.NewFromConfig(cfg)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	queueName := os.Getenv("QUERY_QUEUE")
	catalog, found := request.PathParameters["catalog"]
	if !found {
		log.Fatalln("no catalog variable in path")
	}
	owner, found := request.PathParameters["owner"]
	if !found {
		log.Fatalln("no owner variable in path")
	}
	name, found := request.PathParameters["name"]
	if !found {
		log.Fatalln("no name variable in path")
	}

	client := sqsSession(ctx)
	// client := lambdaSession(ctx)
	// func_name := "queryScore"
	// repoInfo := util.RepoRequestInfo{
	// 	Catalog:   catalog,
	// 	Owner:     owner,
	// 	Name:      name,
	// 	TimeFrame: 6, //temp hardcoded
	// }

	// payload, err := json.Marshal(repoInfo)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// params := lambda.InvokeInput{
	// 	FunctionName:   &func_name,
	// 	Payload:        payload,
	// 	InvocationType: types.InvocationTypeEvent,
	// }

	// _, invoke_err := client.Invoke(ctx, &params)
	// if invoke_err != nil {
	// 	log.Fatalln(err)
	// }
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &queueName,
	}

	result, err := client.GetQueueUrl(ctx, gQInput)
	if err != nil {
		log.Println("Got an error getting the queue URL:")
		log.Println(err)
		return events.APIGatewayProxyResponse{StatusCode: 503, Body: string("Error while getting the queue URL")}, err
	}

	queueURL := result.QueueUrl
	timeFrame := "6"
	sMInput := &sqs.SendMessageInput{
		MessageGroupId: aws.String("handler"),
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
		QueueUrl:    queueURL,
	}

	log.Println(sMInput)
	_, err = client.SendMessage(ctx, sMInput)
	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		return events.APIGatewayProxyResponse{StatusCode: 503, Body: string("Got an error sending the message:")}, err
	}

	message, _ := json.Marshal(response{Message: "Score request accepted"})
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(message)}, nil
}

func main() {
	runtime.Start(handler)
}
