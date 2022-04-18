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

	"golang.org/x/sync/errgroup"
)

const GitUrl = "https://api.github.com/graphql"

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

	postRequest, err := http.NewRequest("POST", GitUrl, responseBody)
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
	repo.CreateDate = data.Data.Repository.CreatedAt
	repo.LatestRelease = data.Data.Repository.LatestRelease.CreatedAt
	repo.Stars = data.Data.Repository.StargazerCount
	repo.DefaultBranch = data.Data.Repository.DefaultBranchRef.Name
	repo.Languages = append(repo.Languages, languages...)

	return nil
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

		post_request, err := http.NewRequest("POST", GitUrl, responseBody)
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
				CreateDate: issueResponse.Created_at,
				Comments:   issueResponse.Comments,
				Assignees:  len(issueResponse.Assignees),
			}
			repo.Issues.OpenIssues = append(repo.Issues.OpenIssues, newIssue)
		} else {
			newIssue := ClosedIssue{
				CreateDate: issueResponse.Created_at,
				CloseDate:  issueResponse.Closed_at,
				Comments:   issueResponse.Comments,
			}
			repo.Issues.ClosedIssues = append(repo.Issues.ClosedIssues, newIssue)
		}
	}

	return len(issues) == 100, nil
}

func GetGithubIssuesRest(client *http.Client, repo *RepoInfo, startDate string) error {
	errs, _ := errgroup.WithContext(context.Background())
	closedHasNextPage := true
	openHasNextPage := true
	closePage := 1
	openPage := 1

	errs.Go(func() error {
		var err error
		for closedHasNextPage {
			closedHasNextPage, err = getGithubIssuePage(client, repo, "closed", closePage, startDate)
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
			openHasNextPage, err = getGithubIssuePage(client, repo, "open", openPage, startDate)
			if err != nil {
				return err
			}
			openPage += 1
		}
		return nil
	})

	return errs.Wait()
}

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
				CreateDate: pr.Created_at,
			}
			repo.PullRequests.OpenPR = append(repo.PullRequests.OpenPR, newPR)
		} else {
			newPR := ClosedPR{
				CreateDate: pr.Created_at,
				CloseDate:  pr.Closed_at,
			}
			repo.PullRequests.ClosedPR = append(repo.PullRequests.ClosedPR, newPR)
		}
	}

	return len(prs) == 100, nil
}

func GetGithubPullRequestsRest(client *http.Client, repo *RepoInfo, startDate string) error {
	errs, _ := errgroup.WithContext(context.Background())
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

func CheckRepoAccess(client *http.Client, owner string, name string) (int, error) {
	requestUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, name)

	body := bytes.NewBuffer(make([]byte, 0))
	request, err := http.NewRequest("GET", requestUrl, body)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return 1, nil
	case 403: // rate limiting exceeded
		return -1, nil
	case 404: // repo not found
		return 0, nil
	default:
		return 0, nil
	}
}

func GetGithubCommitsRest(client *http.Client, repo *RepoInfo, startDate string) error {
	hasNextPage := true
	var err error
	page := 1

	for hasNextPage {
		hasNextPage, err = getGithubCommitsPage(client, repo, page, startDate)
		if err != nil {
			log.Println(err)
			return err
		}
		page += 1
	}

	return nil
}

// deprecated
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

		post_request, err := http.NewRequest("POST", GitUrl, responseBody)
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
					CreateDate:   node.Node.CreatedAt,
					CloseDate:    node.Node.ClosedAt,
					Participants: node.Node.Participants.TotalCount,
					Comments:     node.Node.Assignees.TotalCount,
				}
				closedIssues = append(closedIssues, issue)
			} else {
				issue := OpenIssue{
					CreateDate:   node.Node.CreatedAt,
					Assignees:    node.Node.Assignees.TotalCount,
					Participants: node.Node.Participants.TotalCount,
					Comments:     node.Node.Assignees.TotalCount,
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

		post_request, err := http.NewRequest("POST", GitUrl, responseBody)
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

func GetGithubReleases(client *http.Client, repo *RepoInfo, startDate string) error {
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

		post_request, err := http.NewRequest("POST", GitUrl, responseBody)
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
				CreateDate: node.Node.CreatedAt,
			})
		}
		hasNextPage = data.Data.Repository.Releases.PageInfo.HasNextPage
		cursor = data.Data.Repository.Releases.PageInfo.EndCursor
	}

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
