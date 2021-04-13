package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/jamesjj/qkr"
	"strings"
)

var (
	AcceptableOrigins = map[string]bool{
		"http://localhost:1313": true,
	}
)

func main() {
	runtime.Start(handleRequest)
}

func handleRequest(ctx context.Context, lambdaEvent events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	responseHeaders := map[string]string{
		qkr.Canonical("content-type"): "application/json;charset=utf-8",
	}

	if lambdaEvent.HTTPMethod == "OPTIONS" {
		responseHeaders[qkr.Canonical("access-control-allow-headers")] = strings.Join([]string{
			qkr.Canonical("accept"),
			qkr.Canonical("accept-encoding"),
			qkr.Canonical("content-type"),
		}, ",")
		responseHeaders[qkr.Canonical("access-control-allow-methods")] = "POST"
		responseHeaders[qkr.Canonical("access-control-max-age")] = "3600"
		return events.APIGatewayProxyResponse{Headers: responseHeaders, StatusCode: 200}, nil
	}

	responseBody := map[string]string{
		"foo": "bar",
	}

	return events.APIGatewayProxyResponse{
		IsBase64Encoded: true,
		Body:            qkr.JSONBase64(responseBody),
		Headers:         responseHeaders,
		StatusCode:      200,
	}, nil
}
