package apiutil

import (
	"PennazSoftware/ghin-analyzer/pkg/hcutil"
	"PennazSoftware/ghin-analyzer/services/backend/ghin-api/model"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// GetErrorResponse generates an API Gateway response containing error information
func GetErrorResponse(statusCode int, errorMessage string) (events.APIGatewayProxyResponse, error) {
	errorResponse := model.ErrorResponse{
		ResponseBase: model.ResponseBase{
			Status:       "error",
			ErrorMessage: errorMessage,
			Timestamp:    hcutil.GetCurrentTimestamp(),
		},
	}

	log.Printf("Error (%d): %s", statusCode, errorMessage)

	errorResponseBytes, err := json.Marshal(errorResponse)
	if err != nil {
		log.Printf("failed to marshal errorResponse - %+v", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		Body:       string(errorResponseBytes),
	}, nil
}

// BuildResponse creates the response object
func BuildResponse(response interface{}) events.APIGatewayProxyResponse {
	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("failed to marshal response body: %+v", err)
	}

	statusCode := http.StatusOK

	// If the response object has a status field with a value of "error", set the status code to 500
	if resp, ok := response.(model.ResponseBase); ok && resp.Status == "error" {
		statusCode = http.StatusBadRequest
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			"Access-Control-Allow-Methods": "*",
		},
		Body: string(responseBytes),
	}
}
