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

	// description, err := util.GetRepoLicense(*client, context.Background(), "facebook", "react")
	description, err := util.GetAllIssues(*client, context.Background(), "fonarevvichka", "OSScore")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(description)
}
