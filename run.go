package main

import (
	"context"
	"fmt"
	"go_exploring/util"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
)

func main() {
	gitUrl := "https://api.github.com/graphql"
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_PAT")},
	)
	client := oauth2.NewClient(context.Background(), src)

	repoOwner := "swagger-api"
	repoName := "swagger-ui"

	repoInfo, err := (util.GetCoreRepoInfo(client, gitUrl, repoOwner, repoName))
	if err != nil {
		log.Fatalln(err)
	}

	startDate := time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC).Format(time.RFC3339)

	repoInfo.Issues = util.GetIssues(client, gitUrl, repoOwner, repoName, startDate)
	repoInfo.Dependencies = util.GetDependencies(client, gitUrl, repoOwner, repoName)

	fmt.Println(repoInfo)
}
