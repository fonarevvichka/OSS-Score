package main

import (
	"api/util"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
)

type response struct {
	Message string
	Score   util.Score
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

	mongoClient := util.GetMongoClient()
	score, scoreStatus := util.GetCachedScore(mongoClient, catalog, owner, name, scoreType, 12) // TEMP HARDCODED TO 12 MONTHS

	var message string
	if scoreStatus == 0 {
		message = "Score not yet calculated"
	} else if scoreStatus == 1 {
		message = "Score calculation in progress"
	} else {
		message = "Score ready"
	}
	// retrieve score from database
	//if score not in database send wait / error message
	//if score in database send score

	response, _ := json.Marshal(response{Message: message, Score: score})
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(response)}, nil
}

func main() {
	runtime.Start(handler)
}
