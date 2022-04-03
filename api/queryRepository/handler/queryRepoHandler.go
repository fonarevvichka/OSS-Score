package main

import (
	"api/util"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

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

	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Methods": "POST",
	}

	log.Println("Ready to submit request for ", owner, "/", name)
	// CHECK IF REPO IS VALID AND PUBLIC
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_PAT")},
	)
	httpClient := oauth2.NewClient(ctx, src)

	access, err := util.CheckRepoAccess(httpClient, owner, name)
	if err != nil {
		log.Println(err)
	}

	if access == 0 {
		message, _ := json.Marshal(response{Message: "Could not access repo, check that it was inputted correctly and is public"})
		return events.APIGatewayProxyResponse{
			StatusCode: 406,
			Headers:    headers,
			Body:       string(message),
		}, nil
	} else if access == -1 {
		message, _ := json.Marshal(response{Message: "Github API rate limiting exceeded, cannot submit new repos at this time"})
		return events.APIGatewayProxyResponse{
			StatusCode: 503,
			Headers:    headers,
			Body:       string(message),
		}, nil
	}

	var body util.ScoreRequestBody
	timeFrame := 12 // default to 12 months
	if request.Body != "" {
		err := json.Unmarshal([]byte(request.Body), &body)
		if err != nil {
			message, _ := json.Marshal(response{Message: "Error parsing body of request"})
			return events.APIGatewayProxyResponse{
				StatusCode: 409, //TODO: might need a more accurate code
				Headers:    headers,
				Body:       string(message),
			}, nil
		}
		timeFrame = body.TimeFrame
	}

	client := util.GetSqsClient(ctx)

	queueURL := os.Getenv("QUEUE_URL")
	messageBody := fmt.Sprintf("%s/%s/%d", owner, name, timeFrame)
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
				StringValue: aws.String(strconv.Itoa(timeFrame)),
			},
		},
		MessageBody: aws.String(messageBody),
		QueueUrl:    &queueURL,
	}

	_, err = client.SendMessage(ctx, sMInput)
	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		return events.APIGatewayProxyResponse{StatusCode: 503, Body: string("Got an error sending the message:")}, err
	}

	mongoClient, connected, err := util.GetMongoClient(ctx)
	if connected {
		defer mongoClient.Disconnect(ctx)
	}

	if err != nil {
		message, _ := json.Marshal(response{Message: "Error connecting to MongoDB"})
		return events.APIGatewayProxyResponse{
			StatusCode: 501,
			Headers:    headers,
			Body:       string(message),
		}, err
	}

	collection := mongoClient.Database(os.Getenv("MONGO_DB")).Collection(catalog)

	err = util.SetScoreState(ctx, collection, owner, name, 1)

	if err != nil {
		message, _ := json.Marshal(response{Message: "Error updating state in MongoDB"})
		return events.APIGatewayProxyResponse{
			StatusCode: 501,
			Headers:    headers,
			Body:       string(message),
		}, err
	}

	response, _ := json.Marshal(response{Message: "Score calculation request queued"})
	resp := events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers:    headers,
		Body:       string(response),
	}

	return resp, nil
}

func main() {
	runtime.Start(handler)
}
