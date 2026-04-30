package awsdynamodb

import (
	"PennazSoftware/ghin-analyzer/pkg/hcmodel"
	"PennazSoftware/ghin-analyzer/pkg/hcutil"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// CreateGolfer creates a new Golfer. If successful, it returns the Golfer with its ID
func (c *Client) CreateGolfer(golfer hcmodel.Golfer) (hcmodel.Golfer, error) {

	// Update fields
	timestamp := hcutil.GetCurrentTimestamp()
	golfer.CreatedAt = timestamp
	golfer.UpdatedAt = timestamp

	item, err := attributevalue.MarshalMap(golfer)
	if err != nil {
		log.Printf("Failed to marshal Golfer: %v", err)
		return golfer, err
	}

	input := &dynamodb.PutItemInput{
		TableName: &c.GolfersTableName,
		Item:      item,
	}

	_, err = c.dynamoClient.PutItem(c.Context, input)
	if err != nil {
		log.Printf("Failed to put item in DynamoDB: %v", err)
		return golfer, err
	}

	return golfer, nil
}

// GetGolfer retrieves a Golfer given its ID
func (c *Client) GetGolfer(golferID int) (hcmodel.Golfer, error) {
	var golfer hcmodel.Golfer

	keyEx := expression.Key("golferID").Equal(expression.Value(golferID))

	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return golfer, err
	}

	// Assemble Query
	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(c.GolfersTableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	result, err := c.dynamoClient.Query(c.Context, queryInput)
	if err != nil {
		return golfer, err
	}

	// Since we can technically have multiple results returned, we only care about the first one
	if len(result.Items) > 0 {
		// Deserialize the result
		err = attributevalue.UnmarshalMap(result.Items[0], &golfer)
		if err != nil {
			return golfer, err
		}
	} else {
		return golfer, fmt.Errorf("Golfer not found for golfer (%s)", golferID)
	}

	return golfer, nil
}

// GetAllGolfers gets all Golfers in a paginated manner
func (c *Client) GetAllGolfers(startKey map[string]types.AttributeValue) (golfers []hcmodel.Golfer, pageKey map[string]types.AttributeValue, err error) {
	scanInput := &dynamodb.ScanInput{
		TableName:         aws.String(c.GolfersTableName),
		ExclusiveStartKey: startKey,
	}

	result, err := c.dynamoClient.Scan(c.Context, scanInput)
	if err != nil {
		return golfers, nil, err
	}

	// Deserialize each item in the result
	for _, item := range result.Items {
		var golferScore hcmodel.Golfer
		err = attributevalue.UnmarshalMap(item, &golferScore)
		if err != nil {
			return golfers, nil, err
		}

		golfers = append(golfers, golferScore)
	}

	if result.LastEvaluatedKey != nil {
		// Results are paginated
		pageKey = result.LastEvaluatedKey
	}

	return golfers, pageKey, nil
}

// UpdateGolfer updates a Golfer
func (c *Client) UpdateGolfer(golfer hcmodel.Golfer) (hcmodel.Golfer, error) {

	// Update Timestamp
	golfer.UpdatedAt = hcutil.GetCurrentTimestamp()

	update, err := attributevalue.MarshalMap(golfer)
	if err != nil {
		return golfer, err
	}

	_, err = c.dynamoClient.PutItem(c.Context, &dynamodb.PutItemInput{
		TableName: aws.String(c.GolfersTableName),
		Item:      update,
	})

	if err != nil {
		log.Printf("failed to update Golfer: %+v\nGolfer: %s", err, hcutil.ObjectToJSON(golfer))
	}

	return golfer, err
}

// DeleteGolfer deletes a Golfer given its ID
func (c *Client) DeleteGolfer(golferID int) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(c.GolfersTableName),
		Key: map[string]types.AttributeValue{
			"golferID": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", golferID),
			},
		},
	}

	_, err := c.dynamoClient.DeleteItem(c.Context, input)

	if err != nil {
		log.Printf("failed to delete Golfer for golfer (%d): %+v", golferID, err)
	}

	return err
}
