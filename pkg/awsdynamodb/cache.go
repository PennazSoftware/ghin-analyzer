package awsdynamodb

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// CacheItem represents an item stored in the GHIN API response cache
type CacheItem struct {
	CacheKey string `dynamodbav:"cacheKey"`
	Body     string `dynamodbav:"body"`
	TTL      int64  `dynamodbav:"ttl"`
}

// GetCacheItem retrieves a cached item by its key. Returns nil if not found.
func (c *Client) GetCacheItem(key string) (*CacheItem, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(c.CacheTableName),
		Key: map[string]types.AttributeValue{
			"cacheKey": &types.AttributeValueMemberS{Value: key},
		},
	}

	result, err := c.dynamoClient.GetItem(c.Context, input)
	if err != nil {
		log.Printf("GetCacheItem() error retrieving cache item for key %s: %v", key, err)
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var item CacheItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		log.Printf("GetCacheItem() error unmarshalling cache item for key %s: %v", key, err)
		return nil, err
	}

	return &item, nil
}

// SetCacheItem stores an item in the DynamoDB cache table with the specified Unix timestamp TTL
func (c *Client) SetCacheItem(key string, body string, ttl int64) error {
	item := CacheItem{
		CacheKey: key,
		Body:     body,
		TTL:      ttl,
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		log.Printf("SetCacheItem() error marshalling cache item for key %s: %v", key, err)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(c.CacheTableName),
		Item:      av,
	}

	_, err = c.dynamoClient.PutItem(c.Context, input)
	if err != nil {
		log.Printf("SetCacheItem() error storing cache item for key %s: %v", key, err)
	}

	return err
}
