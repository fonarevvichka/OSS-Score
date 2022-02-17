package main

import (
	"api/util"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/mongo"
)

type response struct {
	Message    string  `json:"message"`
	Metric     float64 `json:"metric"`
	Confidence int     `json:"confidence"`
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
	metric, found := request.PathParameters["metric"]
	if !found {
		log.Fatalln("no metric variable in path")
	}
	fmt.Printf("%s,%s,%s,%s\n", catalog, owner, name, metric)

	mongoClient := util.GetMongoClient()
	defer mongoClient.Disconnect(context.TODO())
	collection := mongoClient.Database("OSS-Score").Collection(catalog) // TODO MAKE DB NAME ENV VAR

	res := util.GetRepoFromDB(collection, owner, name)

	var repo util.RepoInfo
	var metricValue float64
	var confidence int
	var message string

	if res.Err() != mongo.ErrNoDocuments { // match in DB
		err := res.Decode(&repo)

		if err != nil {
			log.Fatalln(err)
		}
		timeFrame := 12
		startPoint := time.Now().AddDate(-(timeFrame / 12), -(timeFrame % 12), 0)

		switch metric {
		case "stars":
			metricValue = float64(repo.Stars)
			confidence = 100
		case "releaseCadence":
			_, metricValue, confidence = util.ParseReleases(repo.Releases, repo.LatestRelease, startPoint)
		case "ageLastRelease":
			metricValue, _, confidence = util.ParseReleases(repo.Releases, repo.LatestRelease, startPoint)
		case "commitCadence":
			metricValue, _, confidence = util.ParseCommits(repo.Commits, startPoint)
		case "contributors":
			_, contributors, _ := util.ParseCommits(repo.Commits, startPoint)
			metricValue = float64(contributors)
			confidence = 100
		case "issueClosureTime":
			metricValue, confidence = util.ParseIssues(repo.Issues, startPoint)
		default:
			message = fmt.Sprintf("Metric querying not yet supported for %s", metric)
		}
	} else {
		message = "Score not available"
	}

	response, _ := json.Marshal(response{Message: message, Metric: metricValue, Confidence: confidence})
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(response)}, nil
}

func main() {
	runtime.Start(handler)
}
