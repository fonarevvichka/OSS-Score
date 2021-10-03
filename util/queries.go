package util

import (
	"context"

	"github.com/shurcooL/githubv4"
)

func GetRepoLicense(client githubv4.Client, ctx context.Context, owner string, name string) (string, error) {
	var q struct {
		Repository struct {
			LicenseInfo struct {
				Key string
				// PsuedoLicense bool
			}
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}

	err := client.Query(ctx, &q, variables)

	return q.Repository.LicenseInfo.Key, err
}

type issue struct {
	Closed bool
	Body   string
	Title  string
	// CreatedAt githubv4.DateTime
	// ClosedAt  githubv4.DateTime
	//comments
}

func GetAllIssues(client githubv4.Client, ctx context.Context, owner string, name string) ([]issue, error) {
	var q struct {
		Repository struct {
			Issues struct {
				Nodes    []issue
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"issues(first: 100, after: $issueCursor)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner":       githubv4.String(owner),
		"name":        githubv4.String(name),
		"issueCursor": (*githubv4.String)(nil),
	}

	var allIssues []issue
	var err error
	for {
		err = client.Query(ctx, &q, variables)
		if err != nil {
			break
		}
		allIssues = append(allIssues, q.Repository.Issues.Nodes...)
		if !q.Repository.Issues.PageInfo.HasNextPage {
			break
		}
		variables["issueCursor"] = githubv4.NewString(q.Repository.Issues.PageInfo.EndCursor)

		if q.Repository.Issues.PageInfo.EndCursor > "400" {
			break
		}
	}
	return allIssues, err
}
