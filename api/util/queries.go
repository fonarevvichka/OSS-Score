package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/google/go-github/v43/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/sync/errgroup"
)

const GraphQLEndpoint = "https://api.github.com/graphql"

func GetCoreRepoInfo(client *http.Client, repo *RepoInfo) error {
	query, err := importQuery("./util/queries/repoInfo.graphql") //TODO: Make this a an env var probably
	if err != nil {
		log.Println(err)
		return err
	}

	variables := fmt.Sprintf("{\"owner\": \"%s\", \"name\": \"%s\"}", repo.Owner, repo.Name)

	postBody, _ := json.Marshal(map[string]string{
		"query":     query,
		"variables": variables,
	})
	responseBody := bytes.NewBuffer(postBody)

	postRequest, err := http.NewRequest("POST", GraphQLEndpoint, responseBody)
	if err != nil {
		log.Println(err)
		return err
	}

	resp, err := client.Do(postRequest)
	if err != nil {
		log.Println(err)
		return err
	}

	if resp.StatusCode != 200 {
		log.Println(resp.Status)
		log.Println(resp.Header)
		log.Println(resp.Body)
		log.Println("Error querying github for core repo info")
		return fmt.Errorf("failed to query github: core repo info\n Status: %s  \n Header: %s \n Body: %s", resp.Status, resp.Header, resp.Body)
	}

	defer resp.Body.Close()

	var data RepoInfoResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		return err
	}

	var languages []string
	for _, node := range data.Data.Repository.Languages.Edges {
		languages = append(languages, node.Node.Name)
	}

	repo.License = data.Data.Repository.LicenseInfo.Key
	repo.CreatedAt = data.Data.Repository.CreatedAt
	repo.LatestRelease = data.Data.Repository.LatestRelease.CreatedAt
	repo.Stars = data.Data.Repository.StargazerCount
	repo.DefaultBranch = data.Data.Repository.DefaultBranchRef.Name
	repo.Languages = append(repo.Languages, languages...)

	return nil
}

func CheckRepoAccess(ctx context.Context, httpClient *http.Client, owner string, name string) (int, error) {
	client := github.NewClient(httpClient)

	// Only need 1 item to check if we can access it
	opts := github.BranchListOptions{
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	}
	_, _, err := client.Repositories.ListBranches(ctx, owner, name, &opts)

	if _, ok := err.(*github.RateLimitError); ok {
		return -1, nil
	} else if _, ok := err.(*github.ErrorResponse); ok { // should probably just check for 404 here not general error
		if err.(*github.ErrorResponse).Message == "Not Found" {
			return 0, nil
		} else {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}

	return 1, nil
}

func GetGithubDependencies(client *http.Client, repo *RepoInfo) error {
	query, err := importQuery("./util/queries/dependencies.graphql")
	if err != nil {
		log.Println(err)
		return err
	}

	var graphCursor string
	var dependencyCursor string

	// hasNextGraphPage := true
	hasNextDependencyPage := true
	// var dependencies []Dependency
	var data DependencyResponse

	// temp: not iterating over all manifests, only primary one
	// for hasNextGraphPage {
	for hasNextDependencyPage {
		variables := fmt.Sprintf("{\"owner\": \"%s\", \"name\": \"%s\", \"graphCursor\": \"%s\", \"dependencyCursor\": \"%s\"}", repo.Owner, repo.Name, graphCursor, dependencyCursor)
		postBody, _ := json.Marshal(map[string]string{
			"query":     query,
			"variables": variables,
		})
		responseBody := bytes.NewBuffer(postBody)

		post_request, err := http.NewRequest("POST", GraphQLEndpoint, responseBody)
		if err != nil {
			log.Println(err)
			return err
		}

		post_request.Header.Add("Accept", "application/vnd.github.hawkgirl-preview+json")
		resp, err := client.Do(post_request)
		if err != nil {
			log.Println(err)
			return err
		}

		if resp.StatusCode != 200 {
			log.Println(resp.Status)
			log.Println(resp.Header)
			log.Println(resp.Body)
			log.Println("Error querying github for dependencies")
			return fmt.Errorf("failed to query github: dependencies\n Status: %s  \n Header: %s \n Body: %s", resp.Status, resp.Header, resp.Body)
		}

		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&data)
		if err != nil {
			log.Println(err)
			return err
		}

		// No dependencies
		if len(data.Data.Repository.DependencyGraphManifests.Edges) == 0 {
			return nil
		}

		for _, node := range data.Data.Repository.DependencyGraphManifests.Edges[0].Node.Dependencies.Edges {
			newDep := Dependency{
				Catalog: "github",
				Owner:   node.Node.Repository.Owner.Login,
				Name:    node.Node.Repository.Name,
				Version: node.Node.Requirements,
			}
			// not pulling enough info out, this shouldn't be needed
			if !dependencyInSlice(newDep, repo.Dependencies) && newDep.Name != "" && newDep.Owner != "" {
				repo.Dependencies = append(repo.Dependencies, newDep)
			}
		}
		hasNextDependencyPage = data.Data.Repository.DependencyGraphManifests.Edges[0].Node.Dependencies.PageInfo.HasNextPage
		dependencyCursor = data.Data.Repository.DependencyGraphManifests.Edges[0].Node.Dependencies.PageInfo.EndCursor
	}
	// hasNextGraphPage = data.Data.Repository.DependencyGraphManifests.PageInfo.HasNextPage
	// graphCursor = data.Data.Repository.DependencyGraphManifests.PageInfo.EndCursor
	// }

	// repo.Dependencies = append(repo.Dependencies, dependencies...)

	return nil
}

func getGithubIssueRestTyped(ctx context.Context, client *github.Client, repo *RepoInfo, opts github.IssueListByRepoOptions) error {
	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, repo.Owner, repo.Name, &opts)

		// May want more granularity here, but for now i think we can check for ratelimitng a step above
		if err != nil {
			return err
		}

		for _, gitIssue := range issues {
			if gitIssue.GetCreatedAt().After(opts.Since) {
				if !gitIssue.IsPullRequest() {
					if *gitIssue.State == "open" { // issue not yet closed
						repo.Issues.OpenIssues = append(repo.Issues.OpenIssues, OpenIssue{
							CreatedAt: gitIssue.GetCreatedAt(),
							Assignees: len(gitIssue.Assignees),
						})
					} else {
						repo.Issues.ClosedIssues = append(repo.Issues.ClosedIssues, ClosedIssue{
							CreatedAt: gitIssue.GetCreatedAt(),
							ClosedAt:  gitIssue.GetClosedAt(),
						})
					}
				}
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return nil
}

func GetGithubIssuesRest(ctx context.Context, httpClient *http.Client, repo *RepoInfo, startPoint time.Time) error {
	errs, ctx := errgroup.WithContext(ctx)
	client := github.NewClient(httpClient)

	opts := github.IssueListByRepoOptions{
		Since: startPoint,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	errs.Go(func() error {
		opts.State = "open"
		return getGithubIssueRestTyped(ctx, client, repo, opts)
	})

	errs.Go(func() error {
		opts.State = "closed"
		return getGithubIssueRestTyped(ctx, client, repo, opts)
	})

	return errs.Wait()
}

func GetGithubCommitsRest(ctx context.Context, httpClient *http.Client, repo *RepoInfo, startPoint time.Time) error {
	client := github.NewClient(httpClient)

	opts := &github.CommitsListOptions{
		Since: startPoint,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		commits, resp, err := client.Repositories.ListCommits(ctx, repo.Owner, repo.Name, opts)

		// May want more granularity here, but for now i think we can check for ratelimitng a step above
		if err != nil {
			return err
		}

		// pull out commits and record them
		for _, gitCommit := range commits {
			repo.Commits = append(repo.Commits, Commit{
				Author:     gitCommit.Commit.Author.GetName(),
				PushedDate: *gitCommit.Commit.Committer.Date,
			})
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return nil
}

func GetGithubReleasesGraphQL(ctx context.Context, httpClient *http.Client, repo *RepoInfo, startPoint time.Time) error {
	client := githubv4.NewClient(httpClient)

	var q struct {
		Repository struct {
			Releases struct {
				Nodes []struct {
					CreatedAt    time.Time
					IsPrerelease bool
				}
				PageInfo PageInfo
			} `graphql:"releases(first: 100, after: $cursor, orderBy: {field: CREATED_AT, direction: DESC})"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":  githubv4.String(repo.Owner),
		"name":   githubv4.String(repo.Name),
		"cursor": (*githubv4.String)(nil), // nil to start n first page
	}

	quit := false
	for {
		err := client.Query(ctx, &q, variables)
		if err != nil {
			return err
		}
		for _, release := range q.Repository.Releases.Nodes {
			if release.CreatedAt.Before(startPoint) {
				quit = true
				break
			}
			if !release.IsPrerelease {
				repo.Releases = append(repo.Releases, Release{
					CreatedAt: release.CreatedAt,
				})
			}
		}

		if !q.Repository.Releases.PageInfo.HasNextPage || quit {
			break
		}
		variables["cursor"] = githubv4.String(q.Repository.Releases.PageInfo.EndCursor)
	}

	// sort releases so most recent release is
	sort.Slice(repo.Releases, func(i, j int) bool {
		return repo.Releases[i].CreatedAt.After(repo.Releases[j].CreatedAt)
	})

	return nil
}

func GetGithubPullsGraphQL(ctx context.Context, httpClient *http.Client, repo *RepoInfo, startPoint time.Time) error {
	client := githubv4.NewClient(httpClient)

	var q struct {
		Repository struct {
			PullRequests struct {
				Nodes []struct {
					Closed    bool
					Merged    bool
					CreatedAt time.Time
					ClosedAt  time.Time
				}
				PageInfo PageInfo
			} `graphql:"pullRequests(first: 100, after: $cursor, orderBy: {field: CREATED_AT, direction: DESC})"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":  githubv4.String(repo.Owner),
		"name":   githubv4.String(repo.Name),
		"cursor": (*githubv4.String)(nil), // nil to start n first page
	}

	quit := false
	for {
		err := client.Query(ctx, &q, variables)
		if err != nil {
			return err
		}

		for _, pull := range q.Repository.PullRequests.Nodes {
			if pull.CreatedAt.Before(startPoint) {
				quit = true
				break
			}
			if pull.Closed {
				if pull.Merged {
					repo.PullRequests.ClosedPR = append(repo.PullRequests.ClosedPR, ClosedPR{
						CreatedAt: pull.CreatedAt,
						ClosedAt:  pull.ClosedAt,
					})
				}
			} else {
				repo.PullRequests.OpenPR = append(repo.PullRequests.OpenPR, OpenPR{
					CreatedAt: pull.CreatedAt,
				})
			}
		}

		if !q.Repository.PullRequests.PageInfo.HasNextPage || quit {
			break
		}
		variables["cursor"] = githubv4.String(q.Repository.PullRequests.PageInfo.EndCursor)
	}

	// sort pull request so most recent release is
	sort.Slice(repo.PullRequests.ClosedPR, func(i, j int) bool {
		return repo.PullRequests.ClosedPR[i].CreatedAt.After(repo.PullRequests.ClosedPR[j].CreatedAt)
	})
	sort.Slice(repo.PullRequests.OpenPR, func(i, j int) bool {
		return repo.PullRequests.OpenPR[i].CreatedAt.After(repo.PullRequests.OpenPR[j].CreatedAt)
	})

	return nil
}

func GetGithubPullsGraphQLManual(client *http.Client, repo *RepoInfo, startPoint time.Time) error {
	query, err := importQuery("./util/queries/pullRequests.graphql") //TODO: Make this a an env var probably
	if err != nil {
		log.Println(err)
		return err
	}

	hasNextPage := true
	stop := false
	cursor := "init"

	var openPulls []OpenPR
	var closedPulls []ClosedPR
	var data PullResponse
	var variables string
	for hasNextPage && !stop {
		if cursor == "init" {
			variables = fmt.Sprintf("{\"owner\": \"%s\", \"name\": \"%s\", \"cursor\": null}", repo.Owner, repo.Name)
		} else {
			variables = fmt.Sprintf("{\"owner\": \"%s\", \"name\": \"%s\", \"cursor\": \"%s\"}", repo.Owner, repo.Name, cursor)
		}
		postBody, _ := json.Marshal(map[string]string{
			"query":     query,
			"variables": variables,
		})
		responseBody := bytes.NewBuffer(postBody)

		post_request, err := http.NewRequest("POST", GraphQLEndpoint, responseBody)
		if err != nil {
			log.Println(err)
			return err
		}

		resp, err := client.Do(post_request)
		if err != nil {
			log.Println(err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Println(resp.Status)
			log.Println(resp.Header)
			log.Println(resp.Body)
			log.Println("Error querying github for issues")
			return fmt.Errorf("failed to query github: issues\n Status: %s  \n Header: %s \n Body: %s", resp.Status, resp.Header, resp.Body)
		}

		decoder := json.NewDecoder(resp.Body)
		if decoder.Decode(&data) != nil {
			log.Println(err)
			return err
		}

		for _, node := range data.Data.Repository.PullRequests.Edges {
			if node.Node.CreatedAt.Before(startPoint) {
				stop = true
				break
			}
			if node.Node.Closed {
				if node.Node.Merged {
					pull := ClosedPR{
						CreatedAt: node.Node.CreatedAt,
						ClosedAt:  node.Node.ClosedAt,
					}
					closedPulls = append(closedPulls, pull)
				}
			} else {
				pull := OpenPR{
					CreatedAt: node.Node.CreatedAt,
				}
				openPulls = append(openPulls, pull)
			}
		}
		hasNextPage = data.Data.Repository.PullRequests.PageInfo.HasNextPage
		cursor = data.Data.Repository.PullRequests.PageInfo.EndCursor
	}

	repo.PullRequests.OpenPR = append(repo.PullRequests.OpenPR, openPulls...)
	repo.PullRequests.ClosedPR = append(repo.PullRequests.ClosedPR, closedPulls...)

	return nil
}

// Takes file path and reads in the query from it
func importQuery(filename string) (string, error) {
	file, err := os.Open(filename)

	if err != nil {
		log.Println(err)
		return "", err
	}

	defer file.Close()

	query, err := ioutil.ReadAll(file)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(query[:]), nil // converts byte array to string
}

// deprecated
func GetGithubReleasesGraphQLManual(client *http.Client, repo *RepoInfo, startDate string) error {
	query, err := importQuery("./util/queries/releases.graphql") //TODO: Make this a an env var probably
	if err != nil {
		log.Println(err)
		return err
	}

	hasNextPage := true
	cursor := "init"

	var data ReleaseResponse
	var variables string

	for hasNextPage {
		if cursor == "init" {
			variables = fmt.Sprintf("{\"owner\": \"%s\", \"name\": \"%s\", \"cursor\": null, \"startDate\": \"%s\"}", repo.Owner, repo.Name, startDate)
		} else {
			variables = fmt.Sprintf("{\"owner\": \"%s\", \"name\": \"%s\", \"cursor\": \"%s\", \"startDate\": \"%s\"}", repo.Owner, repo.Name, cursor, startDate)
		}

		postBody, _ := json.Marshal(map[string]string{
			"query":     query,
			"variables": variables,
		})
		responseBody := bytes.NewBuffer(postBody)

		post_request, err := http.NewRequest("POST", GraphQLEndpoint, responseBody)
		if err != nil {
			log.Println(err)
			return err
		}

		resp, err := client.Do(post_request)
		if err != nil {
			log.Println(err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Println(resp.Status)
			log.Println(resp.Header)
			log.Println(resp.Body)
			log.Println("Error querying github for releases")
			return fmt.Errorf("failed to query github: releases\n Status: %s  \n Header: %s \n Body: %s", resp.Status, resp.Header, resp.Body)
		}

		decoder := json.NewDecoder(resp.Body)
		if decoder.Decode(&data) != nil {
			log.Println(err)
			return err
		}

		for _, node := range data.Data.Repository.Releases.Edges {
			repo.Releases = append(repo.Releases, Release{
				CreatedAt: node.Node.CreatedAt,
			})
		}
		hasNextPage = data.Data.Repository.Releases.PageInfo.HasNextPage
		cursor = data.Data.Repository.Releases.PageInfo.EndCursor
	}

	return nil
}

// deprecated
func getGithubIssuePage(client *http.Client, repo *RepoInfo, state string, page int, startDate string) (bool, error) {
	requestUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?", repo.Owner, repo.Name)
	requestUrlWithParams := requestUrl + fmt.Sprintf("page=%d", page) + fmt.Sprintf("&per_page=%d", 100) + fmt.Sprintf("&since=%s", startDate) + fmt.Sprintf("&state=%s", state)

	responseBody := bytes.NewBuffer(make([]byte, 0))
	request, err := http.NewRequest("GET", requestUrlWithParams, responseBody)
	if err != nil {
		log.Println(err)
		return false, err
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return false, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println(resp.Status)
		log.Println(resp.Header)
		log.Println(resp.Body)
		log.Println("Error querying github for issue page")
		return false, fmt.Errorf("failed to query github: issue page \n Status: %s  \n Header: %s \n Body: %s", resp.Status, resp.Header, resp.Body)
	}

	issues := []IssueResponseRest{}
	decoder := json.NewDecoder(resp.Body)
	if decoder.Decode(&issues) != nil {
		log.Println(err)
		return false, err
	}

	// Pull out issue info
	for _, issueResponse := range issues {
		if issueResponse.State == "open" {
			newIssue := OpenIssue{
				CreatedAt: issueResponse.Created_at,
				// Comments:   issueResponse.Comments,
				Assignees: len(issueResponse.Assignees),
			}
			repo.Issues.OpenIssues = append(repo.Issues.OpenIssues, newIssue)
		} else {
			newIssue := ClosedIssue{
				CreatedAt: issueResponse.Created_at,
				ClosedAt:  issueResponse.Closed_at,
				// Comments:   issueResponse.Comments,
			}
			repo.Issues.ClosedIssues = append(repo.Issues.ClosedIssues, newIssue)
		}
	}

	return len(issues) == 100, nil
}

// deprecated
// startDate filter does not seem to be applied
func GetGithubIssuesGraphQL(client *http.Client, repo *RepoInfo, startDate string) error {
	query, err := importQuery("./util/queries/issues.graphql") //TODO: Make this a an env var probably
	if err != nil {
		log.Println(err)
		return err
	}

	hasNextPage := true
	cursor := "init"

	var closedIssues []ClosedIssue
	var openIssues []OpenIssue
	var data IssueResponse
	var variables string
	for hasNextPage {
		if cursor == "init" {
			variables = fmt.Sprintf("{\"owner\": \"%s\", \"name\": \"%s\", \"cursor\": null, \"startDate\": \"%s\"}", repo.Owner, repo.Name, startDate)
		} else {
			variables = fmt.Sprintf("{\"owner\": \"%s\", \"name\": \"%s\", \"cursor\": \"%s\", \"startDate\": \"%s\"}", repo.Owner, repo.Name, cursor, startDate)
		}
		postBody, _ := json.Marshal(map[string]string{
			"query":     query,
			"variables": variables,
		})
		responseBody := bytes.NewBuffer(postBody)

		post_request, err := http.NewRequest("POST", GraphQLEndpoint, responseBody)
		if err != nil {
			log.Println(err)
			return err
		}

		resp, err := client.Do(post_request)
		if err != nil {
			log.Println(err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Println(resp.Status)
			log.Println(resp.Header)
			log.Println(resp.Body)
			log.Println("Error querying github for issues")
			return fmt.Errorf("failed to query github: issues\n Status: %s  \n Header: %s \n Body: %s", resp.Status, resp.Header, resp.Body)
		}

		decoder := json.NewDecoder(resp.Body)
		if decoder.Decode(&data) != nil {
			log.Println(err)
			return err
		}

		for _, node := range data.Data.Repository.Issues.Edges {
			if node.Node.Closed {
				issue := ClosedIssue{
					CreatedAt: node.Node.CreatedAt,
					ClosedAt:  node.Node.ClosedAt,
					// Participants: node.Node.Participants.TotalCount,
					// Comments:     node.Node.Assignees.TotalCount,
				}
				closedIssues = append(closedIssues, issue)
			} else {
				issue := OpenIssue{
					CreatedAt: node.Node.CreatedAt,
					Assignees: node.Node.Assignees.TotalCount,
					// Participants: node.Node.Participants.TotalCount,
					// Comments:     node.Node.Assignees.TotalCount,
				}
				openIssues = append(openIssues, issue)
			}
		}
		hasNextPage = data.Data.Repository.Issues.PageInfo.HasNextPage
		cursor = data.Data.Repository.Issues.PageInfo.EndCursor
	}

	repo.Issues.OpenIssues = append(repo.Issues.OpenIssues, openIssues...)
	repo.Issues.ClosedIssues = append(repo.Issues.ClosedIssues, closedIssues...)

	return nil
}

// deprecated
// about 70% as fast as using REST
func GetGithubCommitsGraphQL(client *http.Client, repo *RepoInfo, startDate string) error {
	query, err := importQuery("./util/queries/commits.graphql") //TODO: Make this a an env var probably
	if err != nil {
		log.Println(err)
		return err
	}

	hasNextPage := true
	cursor := "init"

	var data CommitResponse
	var variables string

	for hasNextPage {
		if cursor == "init" {
			variables = fmt.Sprintf("{\"owner\": \"%s\", \"name\": \"%s\", \"branch\": \"%s\", \"cursor\": null, \"startDate\": \"%s\"}", repo.Owner, repo.Name, repo.DefaultBranch, startDate)
		} else {
			variables = fmt.Sprintf("{\"owner\": \"%s\", \"name\": \"%s\", \"branch\": \"%s\", \"cursor\": \"%s\", \"startDate\": \"%s\"}", repo.Owner, repo.Name, repo.DefaultBranch, cursor, startDate)
		}

		postBody, _ := json.Marshal(map[string]string{
			"query":     query,
			"variables": variables,
		})
		responseBody := bytes.NewBuffer(postBody)

		post_request, err := http.NewRequest("POST", GraphQLEndpoint, responseBody)
		if err != nil {
			log.Println(err)
			return err
		}

		resp, err := client.Do(post_request)
		if err != nil {
			log.Println(err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Println(resp.Status)
			log.Println(resp.Header)
			log.Println(resp.Body)
			log.Println("Error querying github for commits")
			return fmt.Errorf("failed to query github: commits\n Status: %s  \n Header: %s \n Body: %s", resp.Status, resp.Header, resp.Body)
		}

		decoder := json.NewDecoder(resp.Body)
		if decoder.Decode(&data) != nil {
			log.Println(err)
			return err
		}

		for _, node := range data.Data.Repository.Ref.Target.History.Edges {
			repo.Commits = append(repo.Commits, Commit{
				PushedDate: node.Node.PushedDate,
				Author:     node.Node.Author.Name,
			})
		}
		hasNextPage = data.Data.Repository.Ref.Target.History.PageInfo.HasNextPage
		cursor = data.Data.Repository.Ref.Target.History.PageInfo.EndCursor
	}

	return nil
}

// deprecated
func getGithubCommitsPage(client *http.Client, repo *RepoInfo, page int, startDate string) (bool, error) {
	requestUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?", repo.Owner, repo.Name)
	requestUrlWithParams := requestUrl + fmt.Sprintf("page=%d", page) + fmt.Sprintf("&per_page=%d", 100) + fmt.Sprintf("&since=%s", startDate)

	responseBody := bytes.NewBuffer(make([]byte, 0))
	request, err := http.NewRequest("GET", requestUrlWithParams, responseBody)
	if err != nil {
		log.Println(err)
		return false, err
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return false, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println(resp.Status)
		log.Println(resp.Header)
		log.Println(resp.Body)
		log.Println("Error querying github for core commit page")
		return false, fmt.Errorf("failed to query github: commit page\n Status: %s  \n Header: %s \n Body: %s", resp.Status, resp.Header, resp.Body)
	}

	commits := []CommitResponseRest{}
	decoder := json.NewDecoder(resp.Body)
	if decoder.Decode(&commits) != nil {
		log.Println(err)
		return false, err
	}

	for _, commitResponse := range commits {
		newCommit := Commit{
			Author:     commitResponse.Commit.Author.Name,
			PushedDate: commitResponse.Commit.Author.Date,
		}

		repo.Commits = append(repo.Commits, newCommit)
	}

	return len(commits) == 100, nil
}

// deprecated
func getGithubPullRequestPage(client *http.Client, repo *RepoInfo, state string, page int, startDate string) (bool, error) {
	requestUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?", repo.Owner, repo.Name)
	requestUrlWithParams := requestUrl + fmt.Sprintf("page=%d", page) + fmt.Sprintf("&per_page=%d", 100) + fmt.Sprintf("&since=%s", startDate) + fmt.Sprintf("&state=%s", state)

	responseBody := bytes.NewBuffer(make([]byte, 0))
	request, err := http.NewRequest("GET", requestUrlWithParams, responseBody)
	if err != nil {
		log.Println(err)
		return false, err
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return false, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println(resp.Status)
		log.Println(resp.Header)
		log.Println(resp.Body)
		log.Println("Error querying github for PR page")
		return false, fmt.Errorf("failed to query github: PR page \n Status: %s  \n Header: %s \n Body: %s", resp.Status, resp.Header, resp.Body)
	}

	prs := []PullResponseRest{}
	decoder := json.NewDecoder(resp.Body)
	if decoder.Decode(&prs) != nil {
		log.Println(err)
		return false, err
	}

	// Pull out PR info
	for _, pr := range prs {
		if pr.State == "open" {
			newPR := OpenPR{
				CreatedAt: pr.Created_at,
			}
			repo.PullRequests.OpenPR = append(repo.PullRequests.OpenPR, newPR)
		} else {
			newPR := ClosedPR{
				CreatedAt: pr.Created_at,
				ClosedAt:  pr.Closed_at,
			}
			repo.PullRequests.ClosedPR = append(repo.PullRequests.ClosedPR, newPR)
		}
	}

	return len(prs) == 100, nil
}

// deprecated
func GetGithubPullRequestsRest(ctx context.Context, client *http.Client, repo *RepoInfo, startDate string) error {
	errs, _ := errgroup.WithContext(ctx)
	closedHasNextPage := true
	openHasNextPage := true
	closePage := 1
	openPage := 1

	errs.Go(func() error {
		var err error
		for closedHasNextPage {
			closedHasNextPage, err = getGithubPullRequestPage(client, repo, "closed", closePage, startDate)
			if err != nil {
				return err
			}
			closePage += 1
		}
		return nil
	})

	errs.Go(func() error {
		var err error
		for openHasNextPage {
			openHasNextPage, err = getGithubPullRequestPage(client, repo, "open", openPage, startDate)
			if err != nil {
				return err
			}
			openPage += 1
		}
		return nil
	})

	return errs.Wait()
}
