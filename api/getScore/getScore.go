package main

import (
	"api/util"
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
)

type response struct {
	Message string     `json:"message"`
	Score   util.Score `json:"score"`
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

	dbClient := util.GetDynamoDBSession(ctx)

	// mongoClient := util.GetMongoClient()
	// collection := mongoClient.Database("OSS-Score").Collection(catalog)
	score, scoreStatus := util.GetScore(ctx, dbClient, catalog, owner, name, scoreType, 12) // TEMP HARDCODED TO 12 MONTHS

	var message string
	if scoreStatus == 0 {
		message = "Score not yet calculated"
	} else if scoreStatus == 1 {
		message = "Score calculation queued"
	} else if scoreStatus == 2 {
		message = "Score calculation in progress"
	} else {
		message = "Score ready"
	}
	// retrieve score from database
	//if score not in database send wait / error message
	//if score in database send score

	response, _ := json.Marshal(response{Message: message, Score: score})
	resp := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Methods": "GET",
		},
		Body: string(response)}

	return resp, nil
}

func main() {
	runtime.Start(handler)
}
