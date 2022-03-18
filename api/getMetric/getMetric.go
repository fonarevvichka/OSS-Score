package main

import (
	"api/util"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
)

type singleMetricRepsone struct {
	Message    string  `json:"message"`
	Metric     float64 `json:"metric"`
	Confidence int     `json:"confidence"`
}

type allMetricsResponse struct {
	Message                 string              `json:"message"`
	Stars                   singleMetricRepsone `json:"stars"`
	ReleaseCadence          singleMetricRepsone `json:"releaseCadence"`
	AgeLastRelease          singleMetricRepsone `json:"ageLastRelease"`
	CommitCadence           singleMetricRepsone `json:"commitCadence"`
	Contributors            singleMetricRepsone `json:"contributors"`
	IssueClosureTime        singleMetricRepsone `json:"issueClosureTime"`
	RepoActivityScore       singleMetricRepsone `json:"repoActivityScore"`
	DependencyActivityScore singleMetricRepsone `json:"dependencyActivityScore"`
	RepoLicenseScore        singleMetricRepsone `json:"repoLicenseScore"`
	DependencyLicenseScore  singleMetricRepsone `json:"dependencyLicenseScore"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Methods": "POST",
	}
	
	shelfLife, err := strconv.Atoi(os.Getenv("SHELF_LIFE"))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 501,
			Headers:    headers,
			Body:       "Error converting shelf life env var to int",
		}, err
	}
	
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
	
	timeFrame := 12	
	timeFrameString, found := request.QueryStringParameters["timeFrame"]
	if found {
		var err error
		timeFrame, err = strconv.Atoi(timeFrameString)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 401,
				Headers:    headers,
				Body:       "timeFrame parameter must be an integer",
			}, err
		}
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
	repo, found, err := util.GetRepoFromDB(ctx, collection, owner, name)
	if err != nil {
		log.Fatalln(err)
		//TODO: This should be handeled properly
	}

	var metricValue float64
	var confidence int
	var message string

	var allMetrics allMetricsResponse

	if found { // match in DB
		if repo.Status == 1 {
			message = "Score calculation queued"
		} else if repo.Status == 2 {
			message = "Score calculation in progress"
		} else if repo.Status == 4 {
			message = "Error querying score"
		} else {
			expireDate := time.Now().AddDate(0, 0, -shelfLife)
			startPoint := time.Now().AddDate(-(timeFrame / 12), -(timeFrame % 12), 0)
			
			if repo.UpdatedAt.After(expireDate) && repo.DataStartPoint.Before(startPoint) {
				message = "Metric ready"

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
				case "repoActivityScore":
					score := util.CalculateActivityScore(&repo, startPoint)
					metricValue = score.Score
					confidence = int(score.Confidence)
				case "dependencyActivityScore":
					score, _, _ := util.CalculateDependencyActivityScore(ctx, collection, &repo, startPoint) //TODO: INGORING ERROR
					metricValue = score.Score
					confidence = int(score.Confidence)
				case "repoLicenseScore":
					licenseMap := util.GetLicenseMap()
					score := util.CalculateLicenseScore(&repo, licenseMap)

					metricValue = score.Score
					confidence = int(score.Confidence)
				case "dependencyLicenseScore":
					licenseMap := util.GetLicenseMap()
					score, _, _ := util.CalculateDependencyLicenseScore(ctx, collection, &repo, licenseMap) //TODO: IGNORING ERROR

					metricValue = score.Score
					confidence = int(score.Confidence)
				case "all":
					licenseMap := util.GetLicenseMap()
					var score util.Score

					metricValue = float64(repo.Stars)
					allMetrics.Stars = singleMetricRepsone{
						Metric:     metricValue,
						Confidence: 100,
					}

					_, metricValue, confidence = util.ParseReleases(repo.Releases, repo.LatestRelease, startPoint)
					allMetrics.ReleaseCadence = singleMetricRepsone{
						Metric:     metricValue,
						Confidence: confidence,
					}

					metricValue, _, confidence = util.ParseReleases(repo.Releases, repo.LatestRelease, startPoint)
					allMetrics.AgeLastRelease = singleMetricRepsone{
						Metric:     metricValue,
						Confidence: confidence,
					}

					metricValue, _, confidence = util.ParseCommits(repo.Commits, startPoint)
					allMetrics.CommitCadence = singleMetricRepsone{
						Metric:     metricValue,
						Confidence: confidence,
					}

					_, contributors, _ := util.ParseCommits(repo.Commits, startPoint)
					metricValue = float64(contributors)
					allMetrics.Contributors = singleMetricRepsone{
						Metric:     metricValue,
						Confidence: 100,
					}

					metricValue, confidence = util.ParseIssues(repo.Issues, startPoint)
					allMetrics.IssueClosureTime = singleMetricRepsone{
						Metric:     metricValue,
						Confidence: confidence,
					}

					score = util.CalculateActivityScore(&repo, startPoint)
					metricValue = score.Score
					confidence = int(score.Confidence)
					allMetrics.RepoActivityScore = singleMetricRepsone{
						Metric:     metricValue,
						Confidence: confidence,
					}

					score, _, _ = util.CalculateDependencyActivityScore(ctx, collection, &repo, startPoint) //TODO: INGORING ERROR
					metricValue = score.Score
					confidence = int(score.Confidence)
					allMetrics.DependencyActivityScore = singleMetricRepsone{
						Metric:     metricValue,
						Confidence: confidence,
					}

					score = util.CalculateLicenseScore(&repo, licenseMap)
					metricValue = score.Score
					confidence = int(score.Confidence)
					allMetrics.RepoLicenseScore = singleMetricRepsone{
						Metric:     metricValue,
						Confidence: confidence,
					}

					score, _, _ = util.CalculateDependencyLicenseScore(ctx, collection, &repo, licenseMap) //TODO: IGNORING ERROR
					metricValue = score.Score
					confidence = int(score.Confidence)
					allMetrics.DependencyLicenseScore = singleMetricRepsone{
						Metric:     metricValue,
						Confidence: confidence,
					}
				default:
					message = fmt.Sprintf("Metric querying not yet supported for %s", metric)
				}
			} else {
				message = "Data out of date"
			}
		}
	} else {
		message = "Metric not yet calculated"
	}

	var response []byte

	if metric == "all" {
		allMetrics.Message = message
		response, _ = json.Marshal(allMetrics)
	} else {
		response, _ = json.Marshal(singleMetricRepsone{Message: message, Metric: metricValue, Confidence: confidence})
	}

	resp := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(response),
	}

	return resp, nil
}

func main() {
	runtime.Start(handler)
}
