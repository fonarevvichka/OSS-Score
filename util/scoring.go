package util

import (
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// issues:
// pull out all issues that are < X years old
// of those closed issues calc avg issue closure time
func parseIssues(issues Issues, startPoint time.Time) (float64, int) {
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
func parseCommits(commits []Commit, startPoint time.Time) (float64, int, int) {
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
func parseReleases(releases []Release, LatestRelease time.Time, startPoint time.Time) (float64, float64, int) {
	if len(releases) == 0 {
		// Decrease confidence score
		return 0.0, 0.0, 0
	}

	var releaseCounter float64

	for _, release := range releases {
		if release.CreateDate.After(startPoint) {
			releaseCounter += 1
		}
	}

	// time since start point converted to months
	timeFrame := time.Since(startPoint).Hours() / 24.0 / 30.0

	return time.Since(LatestRelease).Hours() / 24.0 / 7.0, releaseCounter / timeFrame, 0

}

func minMaxScale(min float64, max float64, val float64) float64 {
	return math.Min((val-min)/(max-min), 1)
}

func CalculateActivityScore(mongoClient *mongo.Client, repoInfo *RepoInfo, startPoint time.Time) (Score, Score) {
	// Weights
	commitWeight := 0.25
	contributorWeight := 0.25
	releaseWeight := 0.25
	issueWeight := 0.25

	commitCadenceWeight := 1.0
	// ageLastCommitWeight := 1 - commitCadenceWeight
	ageLastReleaseWeight := 0.5
	releaseCadenceWeight := 1 - ageLastReleaseWeight

	avgIssueClosureTime, issueConfidence := parseIssues(repoInfo.Issues, startPoint)
	commitCadence, contributors, commitConfidence := parseCommits(repoInfo.Commits, startPoint)
	ageLastRelease, releaseCadence, releaseConfidence := parseReleases(repoInfo.Releases, repoInfo.LatestRelease, startPoint)

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
		Score:      100 * score,
		Confidence: confidence,
	}

	score = 0.0
	confidence = 0.0

	collection := mongoClient.Database("OSS-Score").Collection(repoInfo.Catalog) // TODO MAKE DB NAME ENV VAR
	for _, dependency := range repoInfo.Dependencies {
		res := GetRepoFromDB(collection, dependency.Owner, dependency.Name)

		if res.Err() != mongo.ErrNoDocuments { // match in DB
			var depInfo RepoInfo
			err := res.Decode(&depInfo)

			if err != nil {
				log.Fatalln(err)
			}
			score += depInfo.RepoActivityScore.Score
			confidence += depInfo.RepoActivityScore.Confidence
		}
	}
	numDeps := float64(len(repoInfo.Dependencies))

	depScore := Score{
		Score:      100 * score / numDeps,
		Confidence: confidence / numDeps,
	}

	return repoScore, depScore
}
