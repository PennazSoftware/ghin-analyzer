package courses

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

// searchCourses handles searching for courses based on a search term
func searchCourses(request events.APIGatewayProxyRequest, requestID string, ghinClient *ghin.Client, cacheClient awsdynamodb.AwsDynamoDb) (events.APIGatewayProxyResponse, error) {
	log.Println("searchCourses() start")
	defer log.Println("searchCourses() end")

	search := ""
	val, ok := request.QueryStringParameters["search"]
	if ok {
		search = val
	}

	// Validate
	if search == "" || search == "null" {
		return apiutil.GetErrorResponse(400, "a valid search term is required")
	}

	cacheKey := fmt.Sprintf("%s?search=%s", request.Path, search)

	cached, err := cacheClient.GetCacheItem(cacheKey)
	if err == nil && cached != nil && time.Now().Unix() <= cached.TTL {
		var response model.CourseSearchResponse
		if jsonErr := json.Unmarshal([]byte(cached.Body), &response); jsonErr == nil {
			log.Printf("searchCourses() cache hit for key: %s", cacheKey)
			response.Cached = true
			response.RequestID = requestID
			response.Timestamp = hcutil.GetCurrentTimestamp()
			return apiutil.BuildResponse(response), nil
		}
	}

	// Search for courses
	courseSearchResponse, err := ghinClient.SearchCourses(search)
	if err != nil {
		log.Printf("error searching for courses with search term %s: %s", search, err.Error())
		return apiutil.GetErrorResponse(500, "an error occurred while searching for courses")
	}

	response := model.CourseSearchResponse{
		ResponseBase: model.ResponseBase{
			Status:    "success",
			RequestID: requestID,
			Timestamp: hcutil.GetCurrentTimestamp(),
		},
		Courses: courseSearchResponse.Courses,
	}

	if body, marshalErr := json.Marshal(response); marshalErr == nil {
		if setErr := cacheClient.SetCacheItem(cacheKey, string(body), courseCacheTTL()); setErr != nil {
			log.Printf("searchCourses() failed to cache response for key %s: %v", cacheKey, setErr)
		}
	}

	return apiutil.BuildResponse(response), nil
}
