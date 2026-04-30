package courses

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
	case "/ghin/courses":
		switch request.HTTPMethod {
		case "GET":
			r, _ = searchCourses(request, requestID, ghinClient, cacheClient)
		default:
			return apiutil.GetErrorResponse(http.StatusMethodNotAllowed, fmt.Sprintf("the requested method (%s) is not allowed for %s", request.HTTPMethod, resource))
		}
	case "/ghin/courses/{id}":
		switch request.HTTPMethod {
		case "GET":
			r, _ = getCourse(request, requestID, ghinClient, cacheClient)
		default:
			return apiutil.GetErrorResponse(http.StatusMethodNotAllowed, fmt.Sprintf("the requested method (%s) is not allowed for %s", request.HTTPMethod, resource))
		}
	default:
		return apiutil.GetErrorResponse(http.StatusBadRequest, fmt.Sprintf("the request is invalid for %s", resource))
	}

	return r, nil
}

// courseCacheTTL returns a Unix timestamp one week from now
func courseCacheTTL() int64 {
	return time.Now().Add(7 * 24 * time.Hour).Unix()
}
