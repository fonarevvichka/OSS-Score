package util

import (
	"context"
	"math"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// issues:
// pull out all issues that are < X years old
// of those closed issues calc avg issue closure time
func ParseIssues(issues Issues, startPoint time.Time) (float64, int) {
	var totalClosureTime float64
	var issueCounter float64 = 0
	for _, closedIssue := range issues.ClosedIssues {
		if closedIssue.CreateDate.After(startPoint) {
			totalClosureTime += closedIssue.CloseDate.Sub(closedIssue.CreateDate).Hours()
			issueCounter += 1
		}
	}
	if issueCounter == 0 {
		// here we want to add some confidence score stuff
		return 0, 0
	}
	return (totalClosureTime / 24.0) / issueCounter, 100
}

// commits
// pull out all commits that are < X year old
// note the age of the most recent commit
// in that same loop pull out individual contributors
// number commits / 52 = commits per week
//NET: individual contributors and commit cadence, last commit

// Return: commit cadence, contributors, confidence
func ParseCommits(commits []Commit, startPoint time.Time) (float64, int, int) {
	if len(commits) == 0 {
		return 0, 0, 0
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
		return 0, 0, 0
	}

	// time since start point converted to weeks
	timeFrame := time.Since(startPoint).Hours() / 24.0 / 7.0

	return totalCommits / timeFrame, len(contributorMap), 100
}

// releases
// last release we get for free
// pull out releases that are < X year old
// releases / 12 = releases per month
// NET: Last release age and release cadence
func ParseReleases(releases []Release, LatestRelease time.Time, startPoint time.Time) (float64, float64, int) {
	if len(releases) == 0 {
		// Decrease confidence score
		return 10.0, 10.0, 0
	}

	var releaseCounter float64

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
	return math.Min((val-min)/(max-min), 1)
}

func CalculateActivityScore(repoInfo *RepoInfo, startPoint time.Time) Score {
	// Weights
	commitWeight := 0.25
	contributorWeight := 0.25
	releaseWeight := 0.25
	issueWeight := 0.25

	commitCadenceWeight := 1.0
	ageLastReleaseWeight := 0.5
	releaseCadenceWeight := 1 - ageLastReleaseWeight

	avgIssueClosureTime, issueConfidence := ParseIssues(repoInfo.Issues, startPoint)
	commitCadence, contributors, commitConfidence := ParseCommits(repoInfo.Commits, startPoint)
	ageLastRelease, releaseCadence, releaseConfidence := ParseReleases(repoInfo.Releases, repoInfo.LatestRelease, startPoint)

	// NEEDS MORE RESEARCH FOR ACTUAL VALUES
	issueClosureTimeScore := 1 - minMaxScale(0, 176, avgIssueClosureTime)
	commitCadenceScore := minMaxScale(0, 2, commitCadence)
	contributorScore := minMaxScale(0, 10, float64(contributors))
	ageLastReleaseScore := 1 - minMaxScale(0, 26, ageLastRelease)
	releaseCadenceScore := minMaxScale(0, 0.33, releaseCadence)

	// Scores
	commit_score := (commitCadenceWeight * commitCadenceScore) // + (ageLastCommitWeight * ageLastCommit)
	contributer_score := contributorScore
	release_score := (ageLastReleaseWeight * ageLastReleaseScore) + (releaseCadenceWeight * releaseCadenceScore)
	issue_score := issueClosureTimeScore

	score := (commitWeight * commit_score) +
		(contributorWeight * contributer_score) +
		(releaseWeight * release_score) +
		(issueWeight * issue_score)

	confidence := ((contributorWeight + commitWeight) * float64(commitConfidence)) + (issueWeight * float64(issueConfidence)) + (releaseWeight * float64(releaseConfidence))

	repoScore := Score{
		Score:      10 * score,
		Confidence: confidence,
	}

	return repoScore
}

func CalculateDependencyActivityScore(ctx context.Context, collection *mongo.Collection, repoInfo *RepoInfo, startPoint time.Time) (Score, float64, error) {
	if len(repoInfo.Dependencies) == 0 {
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
	for _, dependency := range repoInfo.Dependencies {
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

			individualScore := CalculateActivityScore(&dep, startPoint)
			score += individualScore.Score
			confidence += individualScore.Confidence

			depsWithScores++
		}(dep, startPoint)
	}
	totalDeps := len(repoInfo.Dependencies)

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

func CalculateLicenseScore(repoInfo *RepoInfo, licenseMap map[string]int) Score {
	licenseScore := 0
	confidence := 100

	license := repoInfo.License

	licenseScore = licenseMap[license] / 10

	// Zero confidence if we can't find the license
	if licenseScore == 0 {
		licenseScore = 10
		confidence = 0
	}

	repoScore := Score{
		Score:      float64(licenseScore),
		Confidence: float64(confidence),
	}

	return repoScore
}

func CalculateDependencyLicenseScore(ctx context.Context, collection *mongo.Collection, repoInfo *RepoInfo, licenseMap map[string]int) (Score, float64, error) {
	if len(repoInfo.Dependencies) == 0 {
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
	for _, dependency := range repoInfo.Dependencies {
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

			individualScore := CalculateLicenseScore(&dep, licenseMap)
			score += individualScore.Score
			confidence += individualScore.Confidence

			depsWithScores++
		}(dep)
	}
	totalDeps := len(repoInfo.Dependencies)

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
