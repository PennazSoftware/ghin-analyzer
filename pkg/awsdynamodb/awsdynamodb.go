package awsdynamodb

import (
	"PennazSoftware/ghin-analyzer/pkg/ghin"
	"PennazSoftware/ghin-analyzer/pkg/hcmodel"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	golferTableBaseName string = "handicap-golfers-"
	scoresTableBaseName string = "handicap-scores-"
	cacheTableBaseName  string = "ghin-api-cache-"

	// AwsUsWest2 is the United States #2 region in the West
	AwsUsWest2 string = "us-west-2"
)

// AwsDynamoDb is the interface for interracting with aws DynamoDB
type AwsDynamoDb interface {
	// Golfer
	CreateGolfer(golfer hcmodel.Golfer) (hcmodel.Golfer, error)
	GetGolfer(golferID int) (hcmodel.Golfer, error)
	GetAllGolfers(startKey map[string]types.AttributeValue) (golfers []hcmodel.Golfer, pageKey map[string]types.AttributeValue, err error)
	UpdateGolfer(golfer hcmodel.Golfer) (hcmodel.Golfer, error)
	DeleteGolfer(golferID int) error

	// Scores
	CreateScore(score ghin.Score) (ghin.Score, error)
	GetScoresByGolfer(golferID int, pageKey map[string]types.AttributeValue) ([]ghin.Score, map[string]types.AttributeValue, error)
	GetScoresByGolferPlayedAtRange(golferID int, startDate, endDate string, pageKey map[string]types.AttributeValue) ([]ghin.Score, map[string]types.AttributeValue, error)
	UpdateScore(score ghin.Score) (ghin.Score, error)
	DeleteScore(scoreID string) error

	// Cache
	GetCacheItem(key string) (*CacheItem, error)
	SetCacheItem(key string, body string, ttl int64) error
}

// Client is the struct that embodies the AWS DynamoDB client
type Client struct {
	dynamoClient     *dynamodb.Client
	Context          context.Context
	GolfersTableName string
	ScoresTableName  string
	CacheTableName   string
}

// verifying if the Client struct is indeed implementing the AwsDynamoDb interface
var _ AwsDynamoDb = (*Client)(nil)

// New creates a new AWS DynamoDB client
func New(environment string) *Client {
	if environment != "dev" && environment != "prod" {
		log.Fatalf("Invalid environment: %s", environment)
	}

	ctx := context.TODO()

	// Load the SDK configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"), // Optional: Specify region explicitly
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create the DynamoDB client
	dbClient := dynamodb.NewFromConfig(cfg)

	log.Printf("Creating DynamoDB client for environment: %s", environment)

	return &Client{
		dynamoClient:     dbClient,
		Context:          ctx,
		GolfersTableName: fmt.Sprintf("%s%s", golferTableBaseName, strings.ToLower(environment)),
		ScoresTableName:  fmt.Sprintf("%s%s", scoresTableBaseName, strings.ToLower(environment)),
		CacheTableName:   fmt.Sprintf("%s%s", cacheTableBaseName, strings.ToLower(environment)),
	}
}
