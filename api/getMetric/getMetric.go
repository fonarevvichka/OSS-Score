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
	_, found := request.PathParameters["catalog"]
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

	dbClient := util.GetDynamoDBClient(ctx)
	repo, found, err := util.GetRepoFromDB(ctx, dbClient, owner, name)
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
		} else {
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
			case "repoActivityScore":
				score := util.CalculateActivityScore(&repo, startPoint)
				metricValue = score.Score
				confidence = int(score.Confidence)
			case "dependencyActivityScore":
				score, _ := util.CalculateDependencyActivityScore(ctx, dbClient, &repo, startPoint) //TODO: INGORING ERROR
				metricValue = score.Score
				confidence = int(score.Confidence)
			case "repoLicenseScore":
				licenseMap := util.GetLicenseMap()
				score := util.CalculateLicenseScore(&repo, licenseMap)

				metricValue = score.Score
				confidence = int(score.Confidence)
			case "dependencyLicenseScore":
				licenseMap := util.GetLicenseMap()
				score, _ := util.CalculateDependencyLicenseScore(ctx, dbClient, &repo, licenseMap) //TODO: IGNORING ERROR

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

				score, _ = util.CalculateDependencyActivityScore(ctx, dbClient, &repo, startPoint) //TODO: INGORING ERROR
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

				score, _ = util.CalculateDependencyLicenseScore(ctx, dbClient, &repo, licenseMap) //TODO: IGNORING ERROR
				metricValue = score.Score
				confidence = int(score.Confidence)
				allMetrics.DependencyLicenseScore = singleMetricRepsone{
					Metric:     metricValue,
					Confidence: confidence,
				}
			default:
				message = fmt.Sprintf("Metric querying not yet supported for %s", metric)
			}
		}
	} else {
		message = "Metric not available"
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
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Methods": "GET",
		},
		Body: string(response),
	}

	return resp, nil
}

func main() {
	runtime.Start(handler)
}
