package main

import (
	"api/util"
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/oauth2"
)

type response struct {
	Message  string     `json:"message"`
	DepRatio float64    `json:"depRatio"`
	Score    util.Score `json:"score"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Methods": "POST",
	}

	//TODO RETURN PROPER ERRORS RATHER THAN JUST LOG.FATAL
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

	timeFrame := 12
	timeFrameString, found := request.QueryStringParameters["timeFrame"]
	if found {
		var err error
		timeFrame, err = strconv.Atoi(timeFrameString)
		if err != nil {
			message, _ := json.Marshal(response{Message: "timeFrame parameter must be an integer"})
			return events.APIGatewayProxyResponse{
				StatusCode: 401,
				Headers:    headers,
				Body:       string(message),
			}, nil
		}
	}

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
		}, err
	} else if access == -1 {
		message, _ := json.Marshal(response{Message: "Github API rate limiting exceeded, cannot verify repo access at this time"})
		return events.APIGatewayProxyResponse{
			StatusCode: 503,
			Headers:    headers,
			Body:       string(message),
		}, nil
	}

	mongoClient, connected, err := util.GetMongoClient(ctx)
	if connected {
		defer mongoClient.Disconnect(ctx)
	}

	if err != nil {
		log.Println(err)
		message, _ := json.Marshal(response{Message: "Error connecting to MongoDB"})
		return events.APIGatewayProxyResponse{
			StatusCode: 501,
			Headers:    headers,
			Body:       string(message),
		}, err
	}

	collection := mongoClient.Database(os.Getenv("MONGO_DB")).Collection(catalog)
	score, depRatio, scoreStatus, message, err := util.GetScore(ctx, collection, catalog, owner, name, scoreType, timeFrame)

	if err != nil {
		log.Println(err)
		message, _ := json.Marshal(response{Message: "Error calculating score"})
		return events.APIGatewayProxyResponse{
			StatusCode: 501,
			Headers:    headers,
			Body:       string(message),
		}, err
	}

	if message == "" {
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
	}
	// retrieve score from database
	//if score not in database send wait / error message
	//if score in database send score

	response, _ := json.Marshal(response{Message: message, Score: score, DepRatio: depRatio})
	resp := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(response)}

	return resp, nil
}

func main() {
	runtime.Start(handler)
}
