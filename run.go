package main

import (
	"context"
	"fmt"
	"go_exploring/util"
	"os"

	"golang.org/x/oauth2"
)

func main() {
	gitUrl := "https://api.github.com/graphql"
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GIT_PAT")},
	)
	client := oauth2.NewClient(context.Background(), src)

	// info, err := (util.GetCoreRepoInfo(client, gitUrl, "facebook", "react"))
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(info)
	fmt.Println(util.GetDependencies(client, gitUrl, "facebook", "react"))
}
