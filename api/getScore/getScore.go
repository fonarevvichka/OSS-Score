package main

import (
	"api/util_v2"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
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

	message, _ := json.Marshal(response{Message: "Score not cached"})
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(message)}, nil
}

func init() {
	util_v2.GetMongoClient()
}

func main() {
	runtime.Start(handler)
}
