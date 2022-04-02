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
	"golang.org/x/oauth2"
)

type singleMetricRepsone struct {
	Message    string  `json:"message"`
	Metric     float64 `json:"metric"`
	Confidence float64 `json:"confidence"`
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
	var err error
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
			message, _ := json.Marshal(singleMetricRepsone{Message: "timeFrame parameter must be an integer"})
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
		message, _ := json.Marshal(singleMetricRepsone{Message: "Could not access repo, check that it was inputted correctly and is public"})
		return events.APIGatewayProxyResponse{
			StatusCode: 406,
			Headers:    headers,
			Body:       string(message),
		}, err
	} else if access == -1 {
		message, _ := json.Marshal(singleMetricRepsone{Message: "Github API rate limiting exceeded, cannot verify repo access at this time"})
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
		message, _ := json.Marshal(singleMetricRepsone{Message: "Error connecting to MongoDB"})
		return events.APIGatewayProxyResponse{
			StatusCode: 501,
			Headers:    headers,
			Body:       string(message),
		}, nil
	}

	collection := mongoClient.Database(os.Getenv("MONGO_DB")).Collection(catalog)
	repo, found, err := util.GetRepoFromDB(ctx, collection, catalog, owner, name)
	if err != nil {
		message, _ := json.Marshal(singleMetricRepsone{Message: "Error getting repos from mongoDB"})
		return events.APIGatewayProxyResponse{
			StatusCode: 406,
			Headers:    headers,
			Body:       string(message),
		}, nil
	}

	var metricValue float64
	var confidence float64
	var message string

	var allMetrics allMetricsResponse
	var score util.Score
	var licenseMap map[string]float64

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
					confidence = 100.0
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
					score = util.CalculateActivityScore(&repo, startPoint)
					metricValue = score.Score
					confidence = score.Confidence
				case "dependencyActivityScore":
					score, _, err = util.CalculateDependencyActivityScore(ctx, collection, &repo, startPoint)
					if err != nil {
						message = err.Error()
						break
					}

					metricValue = score.Score
					confidence = score.Confidence
				case "repoLicenseScore":
					licenseMap, err = util.GetLicenseMap("./util/scores/licenseScoring.csv")
					if err != nil {
						message = "Error accessing license scoring file 1st"
						break
					}
					score = util.CalculateLicenseScore(&repo, licenseMap)

					metricValue = score.Score
					confidence = score.Confidence
				case "dependencyLicenseScore":
					licenseMap, err = util.GetLicenseMap("./util/scores/licenseScoring.csv")
					if err != nil {
						message = "Error accessing license scoring file"
						break
					}

					score, _, altErr := util.CalculateDependencyLicenseScore(ctx, collection, &repo, licenseMap)
					if altErr != nil {
						message = "Error accessing license scoring file"
						break
					}

					metricValue = score.Score
					confidence = score.Confidence
				case "all":
					licenseMap, err = util.GetLicenseMap("./util/scores/licenseScoring.csv")
					if err != nil {
						message = "Error accessing license scoring file"
						break
					}
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
					confidence = score.Confidence
					allMetrics.RepoActivityScore = singleMetricRepsone{
						Metric:     metricValue,
						Confidence: confidence,
					}

					score, _, err = util.CalculateDependencyActivityScore(ctx, collection, &repo, startPoint)
					if err != nil {
						message = "Error calculating activity score"
						break
					}
					metricValue = score.Score
					confidence = score.Confidence
					allMetrics.DependencyActivityScore = singleMetricRepsone{
						Metric:     metricValue,
						Confidence: confidence,
					}

					score = util.CalculateLicenseScore(&repo, licenseMap)
					metricValue = score.Score
					confidence = score.Confidence
					allMetrics.RepoLicenseScore = singleMetricRepsone{
						Metric:     metricValue,
						Confidence: confidence,
					}

					score, _, err = util.CalculateDependencyLicenseScore(ctx, collection, &repo, licenseMap)
					if err != nil {
						message = "Error calculating dependency activity score"
						break
					}
					metricValue = score.Score
					confidence = score.Confidence
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

	return resp, err
}

func main() {
	runtime.Start(handler)
}
