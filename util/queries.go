package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func GetCoreRepoInfo(client *http.Client, gitUrl string, owner string, name string) (RepoInfo, error) {
	query := importQuery("./util/queries/repoInfo.graphql") //TODO: Make this a an env var probably
	variables := fmt.Sprintf("{\"owner\": \"%s\", \"name\": \"%s\"}", owner, name)

	postBody, _ := json.Marshal(map[string]string{
		"query":     query,
		"variables": variables,
	})
	responseBody := bytes.NewBuffer(postBody)

	post_request, err := http.NewRequest("POST", gitUrl, responseBody)
	resp, err := client.Do(post_request)

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

	var langauges []string
	for _, node := range data.Data.Repository.Languages.Edges {
		langauges = append(langauges, node.Node.Name)
	}

	info := RepoInfo{
		License:        data.Data.Repository.LicenseInfo.Key,
		CreateDate:     data.Data.Repository.CreatedAt,
		LatestRealease: data.Data.Repository.LatestRelease.CreatedAt,
		Languages:      langauges,
	}
	return info, err
}

func GetDependencies(client *http.Client, gitUrl string, owner string, name string) []Dependency {
	query := importQuery("./util/queries/dependencies.graphql") //TODO: Make this a an env var probably
	var graphCursor string
	var dependencyCursor string

	hasNextGraphPage := true
	hasNextDependencyPage := true
	var dependencies []Dependency
	var data DependencyResponse

	for hasNextGraphPage { //API is always returning false for some reason
		for hasNextDependencyPage {
			variables := fmt.Sprintf("{\"owner\": \"%s\", \"name\": \"%s\", \"graphCursor\": \"%s\", \"dependencyCursor\": \"%s\"}", owner, name, graphCursor, dependencyCursor)
			postBody, _ := json.Marshal(map[string]string{
				"query":     query,
				"variables": variables,
			})
			responseBody := bytes.NewBuffer(postBody)

			post_request, err := http.NewRequest("POST", gitUrl, responseBody)
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

			// No dependencies found
			if data.Data.Repository.DependencyGraphManifests.TotalCount == 0 {
				return nil
			}

			for _, node := range data.Data.Repository.DependencyGraphManifests.Edges[0].Node.Dependencies.Edges {
				newDep := Dependency{
					PacakgeName:   node.Node.PacakgeName,
					NameWithOwner: node.Node.Repository.NameWithOwner,
					Version:       node.Node.Requirements,
				}
				dependencies = append(dependencies, newDep) //TODO: check for dupes
			}
			hasNextDependencyPage = data.Data.Repository.DependencyGraphManifests.Edges[0].Node.Dependencies.PageInfo.HasNextPage
			dependencyCursor = data.Data.Repository.DependencyGraphManifests.Edges[0].Node.Dependencies.PageInfo.EndCursor
		}
		hasNextGraphPage = data.Data.Repository.DependencyGraphManifests.PageInfo.HasNextPage
		graphCursor = data.Data.Repository.DependencyGraphManifests.PageInfo.EndCursor //TODO this is broken for some reason
	}

	return dependencies
}

func GetIssues(client *http.Client, gitUrl string, owner string, name string, startDate string) Issues {
	query := importQuery("./util/queries/issues.graphql") //TODO: Make this a an env var probably

	hasNextPage := true
	cursor := "init"

	var closedIssues []ClosedIssue
	var openIssues []OpenIssue
	var data IssueResponse
	var variables string
	for hasNextPage {
		if cursor == "init" {
			variables = fmt.Sprintf("{\"owner\": \"%s\", \"name\": \"%s\", \"cursor\": null, \"startDate\": \"%s\"}", owner, name, startDate)
		} else {
			variables = fmt.Sprintf("{\"owner\": \"%s\", \"name\": \"%s\", \"cursor\": \"%s\", \"startDate\": \"%s\"}", owner, name, cursor, startDate)
		}
		postBody, _ := json.Marshal(map[string]string{
			"query":     query,
			"variables": variables,
		})
		responseBody := bytes.NewBuffer(postBody)

		post_request, err := http.NewRequest("POST", gitUrl, responseBody)
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

	return Issues{
		OpenIssues:   openIssues,
		ClosedIssues: closedIssues,
	}
}

// Takes file path and reads in the query from it
func importQuery(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	query, err := ioutil.ReadAll(file)

	return string(query[:]) // converts byte array to string
}
