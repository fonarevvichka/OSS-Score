package main

import (
	"api/util"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"golang.org/x/oauth2"
)

type response struct {
	Message string `json:"message"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	queueName := os.Getenv("QUEUE")
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

	log.Println("Ready to submit request for ", owner, "/", name)
	// CHECK IF REPO IS VALID AND PUBLIC
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_PAT")},
	)
	httpClient := oauth2.NewClient(ctx, src)

	valid, err := util.CheckRepoAccess(httpClient, owner, name)
	if err != nil {
		log.Println(err)
	}

	if !valid {
		message, _ := json.Marshal(response{Message: "Could not access repo, check that it was inputted correctly and is public"})
		return events.APIGatewayProxyResponse{
			StatusCode: 406,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Content-Type",
				"Access-Control-Allow-Methods": "POST",
			},
			Body: string(message),
		}, err
	}

	client := util.GetSqsClient(ctx)

	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &queueName,
	}

	result, err := client.GetQueueUrl(ctx, gQInput)
	if err != nil {
		log.Println("Got an error getting the queue URL:")
		log.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: 503,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Content-Type",
				"Access-Control-Allow-Methods": "POST",
			},
			Body: string("Error while getting the queue URL"),
		}, err
	}

	queueURL := result.QueueUrl
	messageBody := fmt.Sprintf("%s/%s", owner, name)
	sMInput := &sqs.SendMessageInput{
		MessageGroupId: aws.String(messageBody),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"catalog": {
				DataType:    aws.String("String"),
				StringValue: aws.String(catalog),
			},
			"owner": {
				DataType:    aws.String("String"),
				StringValue: aws.String(owner),
			},
			"name": {
				DataType:    aws.String("String"),
				StringValue: aws.String(name),
			},
			"timeFrame": {
				DataType:    aws.String("String"),
				StringValue: aws.String("12"), // temp hardcoded
			},
		},
		MessageBody: aws.String(messageBody),
		QueueUrl:    queueURL,
	}

	_, err = client.SendMessage(ctx, sMInput)
	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		return events.APIGatewayProxyResponse{StatusCode: 503, Body: string("Got an error sending the message:")}, err
	}

	dbClient := util.GetDynamoDBClient(ctx)
	err = util.SetScoreState(ctx, dbClient, catalog, owner, name, 1)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 501,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Content-Type",
				"Access-Control-Allow-Methods": "POST",
			},
			Body: "Error updating state in DynamoDB",
		}, err
	}

	response, _ := json.Marshal(response{Message: "Score calculation request queued"})
	resp := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Methods": "POST",
		},
		Body: string(response),
	}

	return resp, nil
}

func main() {
	runtime.Start(handler)
}
