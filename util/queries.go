package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

const GitUrl = "https://api.github.com/graphql"

func GetCoreRepoInfo(client *http.Client, repo *RepoInfo) {
	query := importQuery("./util/queries/repoInfo.graphql") //TODO: Make this a an env var probably
	variables := fmt.Sprintf("{\"owner\": \"%s\", \"name\": \"%s\"}", repo.Owner, repo.Name)

	postBody, _ := json.Marshal(map[string]string{
		"query":     query,
		"variables": variables,
	})
	responseBody := bytes.NewBuffer(postBody)

	postRequest, err := http.NewRequest("POST", GitUrl, responseBody)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(postRequest)

	//TODO: NEED TO CHECK STATUS CODES HERE VERY IMPORTANT
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var data RepoInfoResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatalln(err)
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
}

func GetGithubDependencies(client *http.Client, repo *RepoInfo) {
	query := importQuery("./util/queries/dependencies.graphql") //TODO: Make this a an env var probably
	var graphCursor string
	var dependencyCursor string

	// hasNextGraphPage := true
	hasNextDependencyPage := true
	var dependencies []Dependency
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
			log.Fatalln(err)
		}

		post_request.Header.Add("Accept", "application/vnd.github.hawkgirl-preview+json")
		resp, err := client.Do(post_request)

		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&data)
		if err != nil {
			log.Fatalln(err)
		}

		if data.Data.Repository.DependencyGraphManifests.TotalCount == 0 {
			break
		}

		for _, node := range data.Data.Repository.DependencyGraphManifests.Edges[0].Node.Dependencies.Edges {
			newDep := Dependency{
				Catalog: "github",
				Owner:   node.Node.Repository.Owner.Login,
				Name:    node.Node.Repository.Name,
				Version: node.Node.Requirements,
			}
			// not pulling enough info out, this shouldn't be needed
			if !dependencyInSlice(newDep, dependencies) {
				dependencies = append(dependencies, newDep)
			}
		}
		hasNextDependencyPage = data.Data.Repository.DependencyGraphManifests.Edges[0].Node.Dependencies.PageInfo.HasNextPage
		dependencyCursor = data.Data.Repository.DependencyGraphManifests.Edges[0].Node.Dependencies.PageInfo.EndCursor
	}
	// hasNextGraphPage = data.Data.Repository.DependencyGraphManifests.PageInfo.HasNextPage
	// graphCursor = data.Data.Repository.DependencyGraphManifests.PageInfo.EndCursor
	// }

	repo.Dependencies = append(repo.Dependencies, dependencies...)
}

func getGithubIssuePage(client *http.Client, repo *RepoInfo, state string, page int, startDate string) bool {
	requestUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?", repo.Owner, repo.Name)
	requestUrlWithParams := requestUrl + fmt.Sprintf("page=%d", page) + fmt.Sprintf("&per_page=%d", 100) + fmt.Sprintf("&since=%s", startDate) + fmt.Sprintf("&state=%s", state)

	responseBody := bytes.NewBuffer(make([]byte, 0))
	request, err := http.NewRequest("GET", requestUrlWithParams, responseBody)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(resp.Status)
	defer resp.Body.Close()

	issues := []IssueResponseRest{}
	decoder := json.NewDecoder(resp.Body)
	if decoder.Decode(&issues) != nil {
		log.Fatalln(err)
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

	return len(issues) == 100
}

func GetGithubIssuesRest(client *http.Client, repo *RepoInfo, startDate string) {
	closedHasNextPage := true
	openHasNextPage := true
	closePage := 1
	openPage := 1

	var wg sync.WaitGroup

	wg.Add(1)
	go func(client *http.Client, repo *RepoInfo, state string, page int, startDate string) {
		defer wg.Done()
		for closedHasNextPage {
			closedHasNextPage = getGithubIssuePage(client, repo, state, closePage, startDate)
			closePage += 1
		}
	}(client, repo, "closed", closePage, startDate)

	wg.Add(1)
	go func(client *http.Client, repo *RepoInfo, state string, page int, startDate string) {
		defer wg.Done()
		for openHasNextPage {
			openHasNextPage = getGithubIssuePage(client, repo, "open", openPage, startDate)
			openPage += 1
		}
	}(client, repo, "open", openPage, startDate)

	wg.Wait()
}

func getGithubCommitsPage(client *http.Client, repo *RepoInfo, page int, startDate string) bool {
	requestUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?", repo.Owner, repo.Name)
	requestUrlWithParams := requestUrl + fmt.Sprintf("page=%d", page) + fmt.Sprintf("&per_page=%d", 100) + fmt.Sprintf("&since=%s", startDate)

	responseBody := bytes.NewBuffer(make([]byte, 0))
	request, err := http.NewRequest("GET", requestUrlWithParams, responseBody)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	commits := []CommitResponseRest{}
	decoder := json.NewDecoder(resp.Body)
	if decoder.Decode(&commits) != nil {
		log.Fatalln(err)
	}

	for _, commitResponse := range commits {
		newCommit := Commit{
			Author:     commitResponse.Commit.Author.Name,
			PushedDate: commitResponse.Commit.Author.Date,
		}

		repo.Commits = append(repo.Commits, newCommit)
	}

	return len(commits) == 100
}

func GetGithubCommitsRest(client *http.Client, repo *RepoInfo, startDate string) {
	hasNextPage := true
	page := 1

	for hasNextPage {
		hasNextPage = getGithubCommitsPage(client, repo, page, startDate)
		page += 1
	}
}

// deprecated
func GetGithubIssuesGraphQL(client *http.Client, repo *RepoInfo, startDate string) {
	query := importQuery("./util/queries/issues.graphql") //TODO: Make this a an env var probably

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
			log.Fatalln(err)
		}

		resp, err := client.Do(post_request)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		if decoder.Decode(&data) != nil {
			log.Fatalln(err)
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
}

// deprecated
func GetGithubCommits(client *http.Client, repo *RepoInfo, startDate string) {
	query := importQuery("./util/queries/commits.graphql") //TODO: Make this a an env var probably

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
			log.Fatalln(err)
		}

		resp, err := client.Do(post_request)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		if decoder.Decode(&data) != nil {
			log.Fatalln(err)
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
}

func GetGithubReleases(client *http.Client, repo *RepoInfo, startDate string) {
	query := importQuery("./util/queries/releases.graphql") //TODO: Make this a an env var probably

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
			log.Fatalln(err)
		}

		resp, err := client.Do(post_request)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		if decoder.Decode(&data) != nil {
			log.Fatalln(err)
		}

		for _, node := range data.Data.Repository.Releases.Edges {
			repo.Releases = append(repo.Releases, Release{
				CreateDate: node.Node.CreatedAt,
			})
		}
		hasNextPage = data.Data.Repository.Releases.PageInfo.HasNextPage
		cursor = data.Data.Repository.Releases.PageInfo.EndCursor
	}
}

// Takes file path and reads in the query from it
func importQuery(filename string) string {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	query, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatalln(err)
	}

	return string(query[:]) // converts byte array to string
}
