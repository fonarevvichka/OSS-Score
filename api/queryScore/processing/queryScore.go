package main

import (
	"context"
	"fmt"
	"log"

	runtime "github.com/aws/aws-lambda-go/lambda"

	util "api/util_v2"
)

type response struct {
	Message string
}

func handler(ctx context.Context, repo util.RepoRequestInfo) (response, error) {
	fmt.Printf("%s,%s,%s\n", repo.Catalog, repo.Owner, repo.Name)
	log.Println("Hey hey look at me")

	message := response{Message: "Score processing request accepted"}

	return message, nil
}

func main() {
	runtime.Start(handler)
}
