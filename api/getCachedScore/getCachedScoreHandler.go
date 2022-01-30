package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/session"
)

type response struct {
	Message string
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
	scoreType, found := request.PathParameters["type"]
	if !found {
		log.Fatalln("no scoreType variable in path")
	}
	fmt.Printf("%s,%s,%s,%s\n", catalog, owner, name, scoreType)

	region := os.Getenv("AWS_REGION")
	session, err := session.NewSession(&aws.Config{ // Use aws sdk to connect to dynamoDB
		Region: &region,
	})

	if err != nil {
		log.Fatalln(err)
	}
	lambda.
	svc := invoke.New(session)

	message, _ := json.Marshal(response{Message: "Score not cached"})
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(message)}, nil
}

// func init() {
// 	fmt.Println("Should connect to the DB here")
// }

func main() {
	lambda.Start(handler)
}
