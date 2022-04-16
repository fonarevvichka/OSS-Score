package util

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetActivityScoringData(path string) (map[string]ScoreCategory, error) {
	categoryMap := make(map[string]ScoreCategory)

	data, err := readCsv(path)

	if err != nil {
		return categoryMap, fmt.Errorf("GetActivityScoringData: %v", err)
	}

	for _, vals := range data {
		category := vals[0]

		weight, err := strconv.ParseFloat(vals[1], 64)
		if err != nil {
			return categoryMap, fmt.Errorf("weight: strconv.Parsefloat: %v", err)
		}
		min, err := strconv.ParseFloat(vals[2], 64)

		if err != nil {
			return categoryMap, fmt.Errorf("min: strconv.Parsefloat: %v", err)
		}

		max, err := strconv.ParseFloat(vals[3], 64)
		if err != nil {
			return categoryMap, fmt.Errorf("max:strconv.Parsefloat: %v", err)
		}

		categoryMap[category] = ScoreCategory{
			Weight: weight,
			Min:    min,
			Max:    max,
		}
	}

	return categoryMap, nil
}

func GetLicenseMap(path string) (map[string]float64, error) {
	data, err := readCsv(path)
	licenseMap := make(map[string]float64)

	if err != nil {
		return licenseMap, fmt.Errorf("GetLicenseMap: %v", err)
	}

	for _, vals := range data {
		name := vals[0]
		score, err := strconv.ParseFloat(vals[1], 64)
		if err != nil {
			return licenseMap, fmt.Errorf("strconv.ParseFloat: %v", err)
		}
		licenseMap[name] = score
	}

	return licenseMap, nil
}

// issues:
// pull out all issues that are < X years old
// of those closed issues calc avg issue closure time
func ParseIssues(issues Issues, startPoint time.Time) (float64, float64) {
	totalClosureTime := 0.0
	closedIssueCounter := 0.0

	for _, issue := range issues.ClosedIssues {
		if issue.CreateDate.After(startPoint) {
			totalClosureTime += issue.CloseDate.Sub(issue.CreateDate).Hours()
			closedIssueCounter += 1
		}
	}

	if closedIssueCounter == 0 {
		openIssueCounter := 0.0
		for _, issue := range issues.OpenIssues {
			if issue.CreateDate.After(startPoint) {
				openIssueCounter += 1
				break
			}
		}

		// no open issues over the time frame either
		if openIssueCounter == 0 {
			return 0, 0
		} else {
			return math.MaxFloat64, 100
		}
	}

	return (totalClosureTime / 24.0) / closedIssueCounter, 100
}

// pull requests:
// pull out all prs that are < X years old
// of those closed prs calc avg pr closure time
func ParsePulls(pulls PullRequests, startPoint time.Time) (float64, float64) {
	totalClosureTime := 0.0
	closedPullCounter := 0.0

	for _, pull := range pulls.ClosedPR {
		if pull.CreateDate.After(startPoint) {
			totalClosureTime += pull.CloseDate.Sub(pull.CreateDate).Hours()
			closedPullCounter += 1
		}
	}

	if closedPullCounter == 0 {
		openPullCounter := 0.0
		for _, pull := range pulls.OpenPR {
			if pull.CreateDate.After(startPoint) {
				openPullCounter += 1
				break
			}
		}

		// no open issues over the time frame either
		if openPullCounter == 0 {
			return 0, 0
		} else {
			return math.MaxFloat64, 100
		}
	}

	return (totalClosureTime / 24.0) / closedPullCounter, 100
}

// commits
// pull out all commits that are < X year old
// note the age of the most recent commit
// in that same loop pull out individual contributors
// number commits / 52 = commits per week
//NET: individual contributors and commit cadence, last commit

// Return: commit cadence, contributors, confidence
func ParseCommits(commits []Commit, startPoint time.Time) (float64, float64, int, float64) {
	if len(commits) == 0 {
		return math.MaxFloat64, 0, 0, 100
	}

	var totalCommits float64
	latestCommit := commits[0].PushedDate
	contributorMap := make(map[string]string)

	for _, commit := range commits {

		if commit.PushedDate.After(startPoint) {
			totalCommits += 1
			if commit.PushedDate.After(latestCommit) {
				latestCommit = commit.PushedDate
			}
			_, ok := contributorMap[commit.Author]
			if !ok {
				contributorMap[commit.Author] = commit.Author
			}
		}
	}

	if totalCommits == 0 {
		return math.MaxFloat64, 0, 0, 100
	}

	// time since start point converted to weeks
	timeFrame := time.Since(startPoint).Hours() / 24.0 / 7.0

	return time.Since(latestCommit).Hours() / 24.0, totalCommits / timeFrame, len(contributorMap), 100
}

// releases
// last release we get for free
// pull out releases that are < X year old
// releases / 12 = releases per month
// NET: Last release age and release cadence
func ParseReleases(releases []Release, LatestRelease time.Time, startPoint time.Time) (float64, float64, float64) {
	// no latest release implies no releases in repo --> max scores with zero confidence
	if LatestRelease.IsZero() {
		return 0.0, math.MaxFloat64, 0.0
	}

	// safegaurd for dividing by 0
	if len(releases) == 0 {
		return math.MaxFloat64, 0.0, 100
	}


	releaseCounter := 0.0
	for _, release := range releases {
		if release.CreateDate.After(startPoint) {
			releaseCounter += 1
		}
	}

	// time since start point converted to months
	timeFrame := time.Since(startPoint).Hours() / 24.0 / 30.0

	return time.Since(LatestRelease).Hours() / 24.0 / 7.0, releaseCounter / timeFrame, 100
}

func minMaxScale(min float64, max float64, val float64) float64 {
	return math.Max(math.Min((val-min)/(max-min), 1), 0)
}

// Metric Score, Metric Confidence
func calculateCategoryScore(metric float64, confidence float64, scoreCategory ScoreCategory, inverse bool) (float64, float64) {
	score := minMaxScale(scoreCategory.Min, scoreCategory.Max, metric)

	if inverse {
		score = 1 - score
	}

	return scoreCategory.Weight * score, scoreCategory.Weight * confidence
}

func CalculateRepoActivityScore(repo *RepoInfo, startPoint time.Time) (Score, error) {
	// Scoring Info
	categoryMap, err := GetActivityScoringData("./util/scores/activityScoring.csv")
	if err != nil {
		log.Println(err)
		return Score{}, fmt.Errorf("CalculateRepoActivityScore: %v", err)
	}

	commitCadenceInfo := categoryMap["commitCadence"]
	issueClosureTimeInfo := categoryMap["issueClosureRate"]
	contributorInfo := categoryMap["contributors"]
	ageLastReleaseInfo := categoryMap["ageLastRelease"]
	releaseCadenceInfo := categoryMap["releaseCadence"]

	// Parse data
	avgIssueClosureTime, issueConfidence := ParseIssues(repo.Issues, startPoint)
	_, commitCadence, contributors, commitConfidence := ParseCommits(repo.Commits, startPoint)
	ageLastRelease, releaseCadence, releaseConfidence := ParseReleases(repo.Releases, repo.LatestRelease, startPoint)

	// NEEDS MORE RESEARCH FOR ACTUAL VALUES
	commitCadenceScore, commitCadenceConfidence := calculateCategoryScore(commitCadence, commitConfidence, commitCadenceInfo, false)
	issueClosureTimeScore, issueClosureTimeConfidence := calculateCategoryScore(avgIssueClosureTime, issueConfidence, issueClosureTimeInfo, true)
	contributorScore, contributorConfidence := calculateCategoryScore(float64(contributors), commitConfidence, contributorInfo, false)
	ageLastReleaseScore, ageLastReleaseConfidence := calculateCategoryScore(ageLastRelease, releaseConfidence, ageLastReleaseInfo, true)
	releaseCadenceScore, releaseCadenceConfidence := calculateCategoryScore(releaseCadence, releaseConfidence, releaseCadenceInfo, false)

	score := commitCadenceScore + contributorScore + ageLastReleaseScore + releaseCadenceScore + issueClosureTimeScore
	confidence := commitCadenceConfidence + issueClosureTimeConfidence + contributorConfidence + ageLastReleaseConfidence + releaseCadenceConfidence

	repoScore := Score{
		Score:      10 * score,
		Confidence: confidence,
	}

	return repoScore, nil
}

func CalculateDependencyActivityScore(ctx context.Context, collection *mongo.Collection, repo *RepoInfo, startPoint time.Time) (Score, float64, error) {
	if len(repo.Dependencies) == 0 {
		return Score{
			Score:      10,
			Confidence: 100,
		}, 1, nil
	}

	var wg sync.WaitGroup
	score := 0.0
	confidence := 0.0
	depsWithScores := 0

	var repos []NameOwner
	for _, dependency := range repo.Dependencies {
		repos = append(repos, NameOwner{
			Owner: dependency.Owner,
			Name:  dependency.Name,
		})
	}

	deps, err := GetReposFromDB(ctx, collection, repos)

	if err != nil {
		return Score{ // not sure if score should be 0 or 100 here
			Score:      10,
			Confidence: 0,
		}, 1, err
	}

	for _, dep := range deps {
		wg.Add(1)
		go func(dep RepoInfo, startPoint time.Time) {
			defer wg.Done()

			individualScore, err := CalculateRepoActivityScore(&dep, startPoint)
			if err != nil {
				log.Println(err)
			} else {
				score += individualScore.Score
				confidence += individualScore.Confidence

				depsWithScores++
			}
		}(dep, startPoint)
	}
	totalDeps := len(repo.Dependencies)

	wg.Wait()

	if depsWithScores != 0 {
		score /= float64(depsWithScores)
		confidence /= float64(depsWithScores)
		confidence *= (float64(depsWithScores) / float64(totalDeps))
	} else {
		score = 10
		confidence = 0
	}

	depRatio := float64(depsWithScores) / float64(totalDeps)
	return Score{
		Score:      score,
		Confidence: confidence,
	}, depRatio, nil
}

func CalculateRepoLicenseScore(repo *RepoInfo, licenseMap map[string]float64) Score {
	licenseScore := 0.0
	confidence := 100.0

	license := repo.License

	// Zero confidence if we can't find the license
	if license == "other" {
		licenseScore = 10
		confidence = 0
	} else {
		licenseScore = licenseMap[license]
	}

	repoScore := Score{
		Score:      licenseScore,
		Confidence: confidence,
	}

	return repoScore
}

func CalculateDependencyLicenseScore(ctx context.Context, collection *mongo.Collection, repo *RepoInfo, licenseMap map[string]float64) (Score, float64, error) {
	if len(repo.Dependencies) == 0 {
		return Score{
			Score:      10,
			Confidence: 100,
		}, 1, nil
	}

	var wg sync.WaitGroup
	score := 0.0
	confidence := 0.0
	depsWithScores := 0.0

	var repos []NameOwner
	for _, dependency := range repo.Dependencies {
		repos = append(repos, NameOwner{
			Owner: dependency.Owner,
			Name:  dependency.Name,
		})
	}
	deps, err := GetReposFromDB(ctx, collection, repos)

	if err != nil {
		return Score{ // not sure if score should be 0 or 100 here
			Score:      10,
			Confidence: 0,
		}, 1, err
	}

	for _, dep := range deps {
		wg.Add(1)
		go func(dep RepoInfo) {
			defer wg.Done()

			individualScore := CalculateRepoLicenseScore(&dep, licenseMap)
			score += individualScore.Score
			confidence += individualScore.Confidence

			depsWithScores++
		}(dep)
	}
	totalDeps := float64(len(repo.Dependencies))

	wg.Wait()

	if depsWithScores != 0 {
		score /= float64(depsWithScores)
		confidence /= depsWithScores
		confidence *= (depsWithScores / totalDeps)
	} else {
		score = 10
		confidence = 0
	}

	depRatio := float64(depsWithScores) / float64(totalDeps)
	return Score{
		Score:      score,
		Confidence: confidence,
	}, depRatio, nil
}
