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
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

// getGolferHandicaps handles retrieving a golfer's handicaps based on their GHIN number
func getGolferHandicaps(request events.APIGatewayProxyRequest, requestID string, ghinClient *ghin.Client, cacheClient awsdynamodb.AwsDynamoDb) (events.APIGatewayProxyResponse, error) {
	log.Println("getGolferHandicaps() start")
	defer log.Println("getGolferHandicaps() end")

	ghinNumber := ""
	val, ok := request.PathParameters["id"]
	if ok {
		ghinNumber = val
	}

	var err error
	courseID := 0
	val, ok = request.QueryStringParameters["course-id"]
	if ok {
		courseID, err = strconv.Atoi(val)
		if err != nil {
			log.Printf("error converting course ID to integer: %s", err.Error())
			return apiutil.GetErrorResponse(400, "course ID must be a valid integer")
		}
	}

	// Validate
	if ghinNumber == "" {
		return apiutil.GetErrorResponse(400, "a valid GHIN number is required")
	}

	if courseID == 0 {
		return apiutil.GetErrorResponse(400, "a valid course ID is required")
	}

	cacheKey := fmt.Sprintf("%s?course-id=%d", request.Path, courseID)

	cached, err := cacheClient.GetCacheItem(cacheKey)
	if err == nil && cached != nil && time.Now().Unix() <= cached.TTL {
		var response model.CourseHandicapsResponse
		if jsonErr := json.Unmarshal([]byte(cached.Body), &response); jsonErr == nil {
			log.Printf("getGolferHandicaps() cache hit for key: %s", cacheKey)
			response.Cached = true
			response.RequestID = requestID
			response.Timestamp = hcutil.GetCurrentTimestamp()
			return apiutil.BuildResponse(response), nil
		}
	}

	// get the course handicaps for the golfer
	courseHandicapsResponse, err := ghinClient.GetCourseHandicaps(ghinNumber, courseID)
	if err != nil {
		log.Printf("error retrieving course handicaps for course ID %d: %s", courseID, err.Error())
		return apiutil.GetErrorResponse(500, "an error occurred while retrieving course handicaps")
	}

	response := model.CourseHandicapsResponse{
		ResponseBase: model.ResponseBase{
			Status:    "success",
			RequestID: requestID,
			Timestamp: hcutil.GetCurrentTimestamp(),
		},
		TeeSets: courseHandicapsResponse.TeeSets,
	}

	if body, marshalErr := json.Marshal(response); marshalErr == nil {
		if setErr := cacheClient.SetCacheItem(cacheKey, string(body), golferCacheTTL()); setErr != nil {
			log.Printf("getGolferHandicaps() failed to cache response for key %s: %v", cacheKey, setErr)
		}
	}

	return apiutil.BuildResponse(response), nil
}
