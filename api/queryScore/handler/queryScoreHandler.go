package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	util "api/util"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
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
	func_name := "queryScore"
	repoInfo := util.RepoRequestInfo{
		Catalog:   catalog,
		Owner:     owner,
		Name:      name,
		TimeFrame: 6, //temp hardcoded
	}

	payload, err := json.Marshal(repoInfo)
	if err != nil {
		log.Fatalln(err)
	}

	params := lambda.InvokeInput{
		FunctionName:   &func_name,
		Payload:        payload,
		InvocationType: types.InvocationTypeEvent,
	}

	_, invoke_err := client.Invoke(ctx, &params)
	if invoke_err != nil {
		log.Fatalln(err)
	}

	message, _ := json.Marshal(response{Message: "Score request accepted"})
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(message)}, nil
}

func main() {
	runtime.Start(handler)
}
