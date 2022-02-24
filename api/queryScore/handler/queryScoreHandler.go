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
)

type response struct {
	Message string
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

	client := util.GetSqsSession(ctx)

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

	mongoClient := util.GetMongoClient()
	defer mongoClient.Disconnect(context.TODO())
	collection := mongoClient.Database("OSS-Score").Collection(catalog) // TODO MAKE DB NAME ENV VAR
	util.UpdateScoreState(collection, catalog, owner, name, 1)

	message, _ := json.Marshal(response{Message: "Score calculation request queued"})
	resp := events.APIGatewayProxyResponse{StatusCode: 200, Headers: make(map[string]string), Body: string(message)}
	resp.Headers["Access-Control-Allow-Methods"] = "OPTIONS,POST,GET"
	resp.Headers["Access-Control-Allow-Headers"] = "Content-Type"
	resp.Headers["Access-Control-Allow-Origin"] = "*"

	return resp, nil
}

func main() {
	runtime.Start(handler)
}
