package main

import (
	"PennazSoftware/ghin-analyzer/services/backend/ghin-api/apiutil"
	"PennazSoftware/ghin-analyzer/services/backend/ghin-api/courses"
	"PennazSoftware/ghin-analyzer/services/backend/ghin-api/golfers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (r events.APIGatewayProxyResponse, e error) {
	environment := strings.ToLower(os.Getenv("ENVIRONMENT"))
	log.Printf("Method %s. Path: %s. Resource: %s. Environment: %s", request.HTTPMethod, request.Path, request.Resource, environment)

	if request.HTTPMethod == "OPTIONS" {
		log.Printf("HandleRequest() OPTIONS request")
		return getOptionsResponse(), nil
	}

	parsedPath := parsePath(request.Path)

	if len(parsedPath) == 0 {
		log.Printf("HandleRequest() the request contains no path info (%s)", parsedPath)
		return apiutil.GetErrorResponse(http.StatusBadRequest, fmt.Sprintf("the request contains no path info (%s)", parsedPath))
	}

	// Route to the main API Gateway Resource
	resource := parsedPath[0]
	lc, _ := lambdacontext.FromContext(ctx)

	switch resource {
	case "ghin":
		switch parsedPath[1] {
		case "courses":
			r, _ = courses.ProcessSubRoute(request, environment, lc.AwsRequestID)
		case "golfers":
			r, _ = golfers.ProcessSubRoute(request, environment, lc.AwsRequestID)
		default:
			return apiutil.GetErrorResponse(http.StatusBadRequest, fmt.Sprintf("the requested ghin route is invalid for endpoint '%s'. Full Route: %s", resource, request.Path))
		}

	default:
		return apiutil.GetErrorResponse(http.StatusBadRequest, fmt.Sprintf("the requested route is invalid for endpoint '%s'", resource))
	}

	if r.Headers != nil {
		r.Headers["Access-Control-Allow-Origin"] = "*"
	} else {
		r.Headers = make(map[string]string)
		r.Headers["Access-Control-Allow-Origin"] = "*"
	}
	r.Headers["Access-Control-Allow-Headers"] = "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token"
	r.Headers["Access-Control-Allow-Methods"] = "*"

	log.Printf("HandleRequest() end")

	return r, nil
}

// parsePath parses the elements of the provided route path
func parsePath(routePath string) []string {
	parsedPath := strings.SplitAfter(routePath, "/")

	paths := parsedPath[1:]

	// cleanup by removing forward slashes
	for i := range paths {
		paths[i] = strings.Replace(paths[i], "/", "", -1)
	}

	return paths
}

func main() {
	lambda.Start(handleRequest)
}

func getOptionsResponse() events.APIGatewayProxyResponse {

	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}

	response.Headers = make(map[string]string)
	response.Headers["Access-Control-Allow-Origin"] = "*"
	response.Headers["Access-Control-Allow-Headers"] = "Authorization,Content-Type,X-Amz-Date,X-Amz-Security-Token,X-Api-Key"

	return response
}
