package util

func ParseIssues() {

}

func parseCommits() {

}

func CalculateActivityScore(repoInfo *RepoInfo) {
	// Scale raw data

	// commits
	// pull out all commits that are < X year old
	// note the age of the most recent commit
	// in that same loop pull out individual contributors
	// number commits / 52 = commits per week
	//NET: individual contributors and commit cadence, last commit
	

	// releases
	// last release we get for free
	// pull out releases that are < X year old
	// releases / 12 = releases per month
	// NET: Last release age and release cadence

	// issues:
	// pull out all issues that are < X years old
	// of those closed issues calc avg issue closure time

	// Vars needed:
	// unknowns
	
	// Weights
	int repoWeight = 0.5
	int depWeight = 1 - repoWeight

	int commitWeight = 0.25
	int contributerWeight = 0.25
	int releaseWeight = 0.25
	int issueWeight = 0.25

	int commitCadenceWeight = 0.5
	int ageLastCommitWeight = 1 - commitCadenceWeight
	int ageLastReleaseWeight = 0.5 
	int release_cadence_weight = 1 - ageLastReleaseWeight

	// Scores
	int commit_score = (commitCadenceWeight * commit_cadence) + (ageLastCommitWeight * age_last_commit)
	int contributer_score = num_contributers
	int release_score = (ageLastReleaseWeight * age_last_release) + (release_cadence_weight * release_cadence)
	int issue_score = avg_issue_closure_time

	int repo_score = (commitWeight * commit_score) + 
				(contributerWeight * contributer_score) +
				(releaseWeight * release_score) + 
				(issueWeight * issue_score)

	int dep_score = ((commitWeight * commit_score) + 
				(contributerWeight * contributer_score) +
				(releaseWeight * release_score) + 
				(issueWeight * issue_score)) / (num_dep)


	int activity_score = (repoWeight * repo_score) + (depWeight * dep_score)
}