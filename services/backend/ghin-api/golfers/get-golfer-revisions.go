package golfers

import (
	"PennazSoftware/ghin-analyzer/pkg/awsdynamodb"
	"PennazSoftware/ghin-analyzer/pkg/ghin"
	"PennazSoftware/ghin-analyzer/pkg/hcutil"
	"PennazSoftware/ghin-analyzer/services/backend/ghin-api/apiutil"
	"PennazSoftware/ghin-analyzer/services/backend/ghin-api/model"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

const dateLayout = "2006-01-02"

// getGolferRevisions handles retrieving a golfer's handicap revision history
func getGolferRevisions(request events.APIGatewayProxyRequest, requestID string, ghinClient *ghin.Client, cacheClient awsdynamodb.AwsDynamoDb) (events.APIGatewayProxyResponse, error) {
	log.Println("getGolferRevisions() start")
	defer log.Println("getGolferRevisions() end")

	ghinNumber := ""
	if val, ok := request.PathParameters["id"]; ok {
		ghinNumber = val
	}

	if ghinNumber == "" {
		return apiutil.GetErrorResponse(400, "a valid GHIN number is required")
	}

	startDateStr, ok := request.QueryStringParameters["start-date"]
	if !ok || startDateStr == "" {
		return apiutil.GetErrorResponse(400, "start-date is required (format: YYYY-MM-DD)")
	}

	endDateStr, ok := request.QueryStringParameters["end-date"]
	if !ok || endDateStr == "" {
		return apiutil.GetErrorResponse(400, "end-date is required (format: YYYY-MM-DD)")
	}

	startDate, err := time.Parse(dateLayout, startDateStr)
	if err != nil {
		return apiutil.GetErrorResponse(400, "start-date must be in YYYY-MM-DD format")
	}

	endDate, err := time.Parse(dateLayout, endDateStr)
	if err != nil {
		return apiutil.GetErrorResponse(400, "end-date must be in YYYY-MM-DD format")
	}

	cacheKey := fmt.Sprintf("%s?start-date=%s&end-date=%s", request.Path, startDateStr, endDateStr)

	cached, err := cacheClient.GetCacheItem(cacheKey)
	if err == nil && cached != nil && time.Now().Unix() <= cached.TTL {
		var response model.RevisionsResponse
		if jsonErr := json.Unmarshal([]byte(cached.Body), &response); jsonErr == nil {
			log.Printf("getGolferRevisions() cache hit for key: %s", cacheKey)
			response.Cached = true
			response.RequestID = requestID
			response.Timestamp = hcutil.GetCurrentTimestamp()
			return apiutil.BuildResponse(response), nil
		}
	}

	revisionResponse, err := ghinClient.GetRevisions(ghinNumber, startDate, endDate)
	if err != nil {
		log.Printf("error retrieving revisions for golfer %s: %s", ghinNumber, err.Error())
		return apiutil.GetErrorResponse(500, "an error occurred while retrieving handicap revisions")
	}

	response := model.RevisionsResponse{
		ResponseBase: model.ResponseBase{
			Status:    "success",
			RequestID: requestID,
			Timestamp: hcutil.GetCurrentTimestamp(),
		},
		Revisions: revisionResponse.Revisions,
	}

	if body, marshalErr := json.Marshal(response); marshalErr == nil {
		if setErr := cacheClient.SetCacheItem(cacheKey, string(body), golferCacheTTL()); setErr != nil {
			log.Printf("getGolferRevisions() failed to cache response for key %s: %v", cacheKey, setErr)
		}
	}

	return apiutil.BuildResponse(response), nil
}
