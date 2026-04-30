package courses

import (
	"PennazSoftware/ghin-analyzer/pkg/awsdynamodb"
	"PennazSoftware/ghin-analyzer/pkg/ghin"
	"PennazSoftware/ghin-analyzer/pkg/hcutil"
	"PennazSoftware/ghin-analyzer/services/backend/ghin-api/apiutil"
	"PennazSoftware/ghin-analyzer/services/backend/ghin-api/model"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

// getCourse handles retrieving a course based on a course ID
func getCourse(request events.APIGatewayProxyRequest, requestID string, ghinClient *ghin.Client, cacheClient awsdynamodb.AwsDynamoDb) (events.APIGatewayProxyResponse, error) {
	log.Println("getCourse() start")
	defer log.Println("getCourse() end")

	var err error
	courseID := 0
	val, ok := request.PathParameters["id"]
	if ok {
		courseID, err = strconv.Atoi(val)
		if err != nil {
			log.Printf("error converting course ID %s to int: %s", val, err.Error())
			return apiutil.GetErrorResponse(400, "a valid course ID is required")
		}
	}

	// Validate
	if courseID == 0 {
		return apiutil.GetErrorResponse(400, "a valid course ID is required")
	}

	cacheKey := request.Path

	cached, err := cacheClient.GetCacheItem(cacheKey)
	if err == nil && cached != nil && time.Now().Unix() <= cached.TTL {
		var response model.CourseDetailsResponse
		if jsonErr := json.Unmarshal([]byte(cached.Body), &response); jsonErr == nil {
			log.Printf("getCourse() cache hit for key: %s", cacheKey)
			response.Cached = true
			response.RequestID = requestID
			response.Timestamp = hcutil.GetCurrentTimestamp()
			return apiutil.BuildResponse(response), nil
		}
	}

	// Fetch course details
	courseSearchResponse, err := ghinClient.GetCourseDetails(courseID)
	if err != nil {
		log.Printf("error retrieving course details for course ID %d: %s", courseID, err.Error())
		return apiutil.GetErrorResponse(500, "an error occurred while retrieving course details")
	}

	response := model.CourseDetailsResponse{
		ResponseBase: model.ResponseBase{
			Status:    "success",
			RequestID: requestID,
			Timestamp: hcutil.GetCurrentTimestamp(),
		},
		Course: courseSearchResponse,
	}

	if body, marshalErr := json.Marshal(response); marshalErr == nil {
		if setErr := cacheClient.SetCacheItem(cacheKey, string(body), courseCacheTTL()); setErr != nil {
			log.Printf("getCourse() failed to cache response for key %s: %v", cacheKey, setErr)
		}
	}

	return apiutil.BuildResponse(response), nil
}
