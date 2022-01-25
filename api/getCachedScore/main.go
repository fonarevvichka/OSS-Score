package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type lambdaEvent struct {
	PathParameters struct {
		Catalog string
		Owner   string
		Name    string
		Type    string
	}
}

type response struct {
	Message string
}

func HandleLambdaEvent(event lambdaEvent) (response, error) {
	catalog := event.PathParameters.Catalog
	owner := event.PathParameters.Owner
	name := event.PathParameters.Name
	Type := event.PathParameters.Type
	return response{Message: fmt.Sprintf("%s,%s,%s,%s", catalog, owner, name, Type)}, nil
}

// func init() {
// 	fmt.Println("Should connect to the DB here")
// }

func main() {
	lambda.Start(HandleLambdaEvent)
}
