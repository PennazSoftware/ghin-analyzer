package golfers

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

// getGolfer handles retrieving a golfer's information based on their GHIN number
func getGolfer(request events.APIGatewayProxyRequest, requestID string, ghinClient *ghin.Client, cacheClient awsdynamodb.AwsDynamoDb) (events.APIGatewayProxyResponse, error) {
	log.Println("getGolfer() start")
	defer log.Println("getGolfer() end")

	var err error
	ghinNumber := 0
	val, ok := request.PathParameters["id"]
	if ok {
		ghinNumber, err = strconv.Atoi(val)
		if err != nil {
			log.Printf("error converting GHIN number to integer: %s", err.Error())
			return apiutil.GetErrorResponse(400, "a valid GHIN number is required")
		}
	}

	// Validate
	if ghinNumber == 0 {
		return apiutil.GetErrorResponse(400, "a valid GHIN number is required")
	}

	cacheKey := request.Path

	cached, err := cacheClient.GetCacheItem(cacheKey)
	if err == nil && cached != nil && time.Now().Unix() <= cached.TTL {
		var response model.GolferResponse
		if jsonErr := json.Unmarshal([]byte(cached.Body), &response); jsonErr == nil {
			log.Printf("getGolfer() cache hit for key: %s", cacheKey)
			response.Cached = true
			response.RequestID = requestID
			response.Timestamp = hcutil.GetCurrentTimestamp()
			return apiutil.BuildResponse(response), nil
		}
	}

	// get the golfer information
	golferResponse, err := ghinClient.GetGolfer(ghinNumber)
	if err != nil {
		log.Printf("error retrieving golfer information: %s", err.Error())
		return apiutil.GetErrorResponse(500, "error retrieving golfer information")
	}

	if len(golferResponse) == 0 {
		return apiutil.GetErrorResponse(404, "golfer not found")
	}

	response := model.GolferResponse{
		ResponseBase: model.ResponseBase{
			Status:    "success",
			RequestID: requestID,
			Timestamp: hcutil.GetCurrentTimestamp(),
		},
		Golfer: golferResponse[0],
	}

	if body, marshalErr := json.Marshal(response); marshalErr == nil {
		if setErr := cacheClient.SetCacheItem(cacheKey, string(body), golferCacheTTL()); setErr != nil {
			log.Printf("getGolfer() failed to cache response for key %s: %v", cacheKey, setErr)
		}
	}

	return apiutil.BuildResponse(response), nil
}
