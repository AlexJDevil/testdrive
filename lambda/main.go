package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/arun0009/testdrive/testdrive"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	switch {

	case strings.ToUpper(request.HTTPMethod) == "GET" && request.Path == "/health":
		log.Println("health endpoint invoked")
		return events.APIGatewayProxyResponse{
			StatusCode:      http.StatusOK,
			Body:            `{"msg": "healthy"}`,
			IsBase64Encoded: false,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil

	default:
		log.Println("testdrive invoked")
		testdrive.Exec()
		return events.APIGatewayProxyResponse{
			StatusCode:      http.StatusOK,
			Body:            `{"msg": "testdrive completed successfully"}`,
			IsBase64Encoded: false,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}
}

func main() {
	lambda.Start(Handler)
}
