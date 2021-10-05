package main

import (
	"context"
	"fmt"
	"go_exploring/util"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_PAT")},
	)

	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	// response, err := util.GetRepoInfo(*client, context.Background(), "facebook", "react")
	response, err := util.GetDependencies(*client, context.Background(), "facebook", "react")
	// description, err := util.GetIssuesByState(*client, context.Background(), "fonarevvichka", "OSScore", githubv4.IssueStateOpen)
	// description, err := util.GetIssuesByState(*client, context.Background(), "facebook", "react", githubv4.IssueStateClosed)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(response)
}
