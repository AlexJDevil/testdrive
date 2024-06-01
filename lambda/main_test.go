package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			request: events.APIGatewayProxyRequest{HTTPMethod: "GET", Path: "/health"},
			expect:  `{"msg": "healthy"}`,
			err:     nil,
		},
	}

	for _, test := range tests {
		response, err := Handler(context.Background(), test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}
