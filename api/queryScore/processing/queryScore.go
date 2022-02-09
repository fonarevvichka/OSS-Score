package main

import (
	"context"
	"fmt"

	runtime "github.com/aws/aws-lambda-go/lambda"

	"api/util"
)

func handler(ctx context.Context, repo util.RepoRequestInfo) error {
	fmt.Printf("%s,%s,%s,%d\n", repo.Catalog, repo.Owner, repo.Name, repo.TimeFrame)

	util.QueryProject(repo.Catalog, repo.Owner, repo.Name, repo.TimeFrame)
	return nil
}

func main() {
	runtime.Start(handler)
}
