package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	util "go_exploring/util"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

func main() {
	gitUrl := "https://api.github.com/graphql"
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_PAT")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	postBody, _ := json.Marshal(map[string]string{
		"query": "query ($name: String!, $owner: String!){	repository(owner: $owner, name: $name) {    name    url  }}",
		"variables": "{\"owner\": \"fonarevvichka\", \"name\": \"go_exploring\"}",
	})
	responseBody := bytes.NewBuffer(postBody)

	post_request, err := http.NewRequest("POST", gitUrl, responseBody)
	post_request.Header.Add("Accept", "application/vnd.github.hawkgirl-preview+json")
	resp, err := httpClient.Do(post_request)

	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var data util.Data
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(data.Data.Repository.Url)
}
