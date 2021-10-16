package main

import (
	"context"
	"fmt"
	util "go_exploring/util"
	"log"
	"os"

	"golang.org/x/oauth2"
)

func main() {
	gitUrl := "https://api.github.com/graphql"
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_PAT")},
	)
	client := oauth2.NewClient(context.Background(), src)

	// var repoInfo util.RepoInfo
	info, err := (util.GetRepoInfo(client, gitUrl, "facebook", "react"))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(info)

}
