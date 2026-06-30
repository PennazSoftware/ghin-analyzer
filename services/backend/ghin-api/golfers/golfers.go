package golfers

import (
	"PennazSoftware/ghin-analyzer/pkg/awsdynamodb"
	"PennazSoftware/ghin-analyzer/pkg/awssecretsmanager"
	"PennazSoftware/ghin-analyzer/pkg/ghin"
	"PennazSoftware/ghin-analyzer/services/backend/ghin-api/apiutil"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

const (
	ghinSecretName = "Ghin_50155"
)

// ProcessSubRoute handles processing of the subroute resources
func ProcessSubRoute(request events.APIGatewayProxyRequest, environment string, requestID string) (r events.APIGatewayProxyResponse, e error) {
	resource := request.Resource

	// Get the GHIN credentials from AWS Secrets Manager
	secretsClient := awssecretsmanager.New()
	var ghinSecret awssecretsmanager.GHINSecret
	err := secretsClient.GetSecret(ghinSecretName, &ghinSecret)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ghinClient := ghin.NewWithLogin(ghinSecret.Username, ghinSecret.Password)
	cacheClient := awsdynamodb.New(environment)

	switch resource {
	case "/ghin/golfers/{id}":
		switch request.HTTPMethod {
		case "GET":
			r, _ = getGolfer(request, requestID, ghinClient, cacheClient)
		default:
			return apiutil.GetErrorResponse(http.StatusMethodNotAllowed, fmt.Sprintf("the requested method (%s) is not allowed for %s", request.HTTPMethod, resource))
		}
	case "/ghin/golfers/{id}/handicaps":
		switch request.HTTPMethod {
		case "GET":
			r, _ = getGolferHandicaps(request, requestID, ghinClient, cacheClient)
		default:
			return apiutil.GetErrorResponse(http.StatusMethodNotAllowed, fmt.Sprintf("the requested method (%s) is not allowed for %s", request.HTTPMethod, resource))
		}
	case "/ghin/golfers/{id}/revisions":
		switch request.HTTPMethod {
		case "GET":
			r, _ = getGolferRevisions(request, requestID, ghinClient, cacheClient)
		default:
			return apiutil.GetErrorResponse(http.StatusMethodNotAllowed, fmt.Sprintf("the requested method (%s) is not allowed for %s", request.HTTPMethod, resource))
		}
	default:
		return apiutil.GetErrorResponse(http.StatusBadRequest, fmt.Sprintf("the request is invalid for %s", resource))
	}

	return r, nil
}

// golferCacheTTL returns a Unix timestamp for midnight of the current day in America/Los_Angeles
func golferCacheTTL() int64 {
	loc, _ := time.LoadLocation("America/Los_Angeles")
	now := time.Now().In(loc)
	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc)
	return midnight.Unix()
}
