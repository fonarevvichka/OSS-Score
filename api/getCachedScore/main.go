package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

type lambdaEvent struct {
}

type response struct {
	Message string
}

func HandleLambdaEvent(event lambdaEvent) (response, error) {
	return response{Message: "success"}, nil
}

// func init() {
// 	fmt.Println("Should connect to the DB here")
// }

func main() {
	lambda.Start(HandleLambdaEvent)
}
