package main

import (
	"context"

	runtime "github.com/aws/aws-lambda-go/lambda"

	"api/util"
)

func handler(ctx context.Context, repo util.RepoRequestInfo) error {
	util.QueryProject(repo.Catalog, repo.Owner, repo.Name, repo.TimeFrame)
	return nil
}

func main() {
	runtime.Start(handler)
}
