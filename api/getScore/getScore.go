package main

import (
	"api/util"
	"context"
	"encoding/json"
	"log"
	"os"

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

	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Methods": "POST",
	}

	mongoClient, connected, err := util.GetMongoClient(ctx)
	if connected {
		defer mongoClient.Disconnect(ctx)
	}

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 501,
			Headers:    headers,
			Body:       "Error connecting to MongoDB",
		}, err
	}

	collection := mongoClient.Database(os.Getenv("MONGO_DB")).Collection(catalog)

	score, scoreStatus := util.GetScore(ctx, collection, catalog, owner, name, scoreType, 12) // TEMP HARDCODED TO 12 MONTHS

	var message string
	if scoreStatus == 0 {
		message = "Score not yet calculated"
	} else if scoreStatus == 1 {
		message = "Score calculation queued"
	} else if scoreStatus == 2 {
		message = "Score calculation in progress"
	} else if scoreStatus == 3 {
		message = "Score ready"
	} else {
		message = "Error querying score"
	}
	// retrieve score from database
	//if score not in database send wait / error message
	//if score in database send score

	response, _ := json.Marshal(response{Message: message, Score: score})
	resp := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(response)}

	return resp, nil
}

func main() {
	runtime.Start(handler)
}
