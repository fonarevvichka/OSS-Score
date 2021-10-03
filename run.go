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

	description, err := getRepoLicense(*client, context.Background(), "facebook", "react")

	if err != nil {
		fmt.Printf("Error! %f \n", err)
		return
	}

	fmt.Println(description)
}

func getRepoLicense(client githubv4.Client, ctx context.Context, owner string, name string) (string, error) {
	var q struct {
		// Repository struct {
		Repository struct {
			Description string
			LicenseInfo struct {
				name string
			}
			// LicenseInfo struct {
			// body string
			// License struct {
			// 	name string
			// }
			// }
			// }
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}
	err := client.Query(ctx, &q, variables)
	// fmt.Println(q.Repository.License)
	return q.Repository.Description, err
}
