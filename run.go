package main

import (
	"context"
	"fmt"
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

	description, err := fetchRepoInfo(*client, context.Background(), "fonarevvichka", "tryNotToLaughAffectiva")

	if err != nil {
		fmt.Println("Error!")
	}

	fmt.Println(description)
}

func fetchRepoInfo(client githubv4.Client, ctx context.Context, owner string, name string) (string, error) {
	var q struct {
		Repository struct {
			Description string
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}

	err := client.Query(ctx, &q, variables)
	return q.Repository.Description, err
}
