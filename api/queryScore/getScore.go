package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

type response struct {
	Message string
}

func lambdaSession(ctx context.Context) *lambda.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	return lambda.NewFromConfig(cfg)
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
	fmt.Printf("%s,%s,%s\n", catalog, owner, name)

	client := lambdaSession(ctx)
	func_name := "dummy"
	params := lambda.InvokeInput{
		FunctionName: &func_name,
	}

	_, err := client.Invoke(ctx, &params)
	if err != nil {
		log.Fatalln(err)
	}

	message, _ := json.Marshal(response{Message: "Score not cached"})
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(message)}, nil
}

// func init() {
// 	fmt.Println("Should connect to the DB here")
// }

func main() {
	runtime.Start(handler)
}
