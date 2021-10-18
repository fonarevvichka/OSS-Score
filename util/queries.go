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

// type issue struct {
// 	Title     string
// 	CreatedAt githubv4.DateTime
// 	ClosedAt  githubv4.DateTime
// }

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
		License:    data.Data.Repository.LicenseInfo.Key,
		CreateDate: data.Data.Repository.CreatedAt,
		Languages:  langauges,
	}
	return info, err
}

func GetDependencies(client *http.Client, gitUrl string, owner string, name string) []Dependency {
	query := importQuery("./util/queries/dependencies.graphql") //TODO: Make this a an env var probably
	graphCursor := ""
	dependencyCursor := ""
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

// func GetIssuesByState(client githubv4.Client, ctx context.Context, owner string, name string, state githubv4.IssueState) ([]issue, error) {
// 	var q struct {
// 		Repository struct {
// 			Issues struct {
// 				Nodes    []issue
// 				PageInfo struct {
// 					EndCursor   githubv4.String
// 					HasNextPage bool
// 				}
// 			} `graphql:"issues(first: 100, after: $issueCursor, states: $states)"`
// 		} `graphql:"repository(owner: $owner, name: $name)"`
// 	}
// 	variables := map[string]interface{}{
// 		"owner":       githubv4.String(owner),
// 		"name":        githubv4.String(name),
// 		"issueCursor": (*githubv4.String)(nil),
// 		"states":      []githubv4.IssueState{state},
// 	}

// 	var allIssues []issue
// 	var err error
// 	for {
// 		err = client.Query(ctx, &q, variables)
// 		if err != nil {
// 			break
// 		}
// 		allIssues = append(allIssues, q.Repository.Issues.Nodes...)
// 		if !q.Repository.Issues.PageInfo.HasNextPage {
// 			break
// 		}
// 		variables["issueCursor"] = githubv4.NewString(q.Repository.Issues.PageInfo.EndCursor)

// 		// if q.Repository.Issues.PageInfo.EndCursor > "400" { // temp to make things quicker
// 		// 	break
// 		// }
// 	}
// 	return allIssues, err
// }
