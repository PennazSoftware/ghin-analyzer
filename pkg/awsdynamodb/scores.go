package awsdynamodb

import (
	"PennazSoftware/ghin-analyzer/pkg/ghin"
	"PennazSoftware/ghin-analyzer/pkg/hcutil"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// CreateScore creates a new Score. If successful, it returns the Score with its ID
func (c *Client) CreateScore(score ghin.Score) (ghin.Score, error) {

	item, err := attributevalue.MarshalMap(score)
	if err != nil {
		log.Printf("Failed to marshal Score: %v", err)
		return score, err
	}

	input := &dynamodb.PutItemInput{
		TableName: &c.ScoresTableName,
		Item:      item,
	}

	_, err = c.dynamoClient.PutItem(c.Context, input)
	if err != nil {
		log.Printf("Failed to put item in DynamoDB: %v", err)
		return score, err
	}

	return score, nil
}

// GetScoresByGolfer retrieves Scores for a given Golfer ID
func (c *Client) GetScoresByGolfer(golferID int, pageKey map[string]types.AttributeValue) ([]ghin.Score, map[string]types.AttributeValue, error) {
	var scores []ghin.Score

	keyEx := expression.Key("golfer_id").Equal(expression.Value(golferID))

	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return scores, nil, err
	}

	// Assemble Query
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(c.ScoresTableName),
		KeyConditionExpression: expr.KeyCondition(),
	}

	if pageKey != nil {
		queryInput.ExclusiveStartKey = pageKey
	}

	result, err := c.dynamoClient.Query(c.Context, queryInput)
	if err != nil {
		return scores, nil, err
	}

	// Deserialize each item in the result
	for _, item := range result.Items {
		var score ghin.Score
		err = attributevalue.UnmarshalMap(item, &score)
		if err != nil {
			return scores, nil, err
		}

		scores = append(scores, score)
	}

	if result.LastEvaluatedKey != nil {
		// Results are paginated
		pageKey = result.LastEvaluatedKey
	}

	return scores, pageKey, nil
}

// GetScoresByGolferPlayedAtRange retrieves Scores for a given Golfer ID within a specific date range
func (c *Client) GetScoresByGolferPlayedAtRange(golferID int, startDate, endDate string, pageKey map[string]types.AttributeValue) ([]ghin.Score, map[string]types.AttributeValue, error) {
	var scores []ghin.Score

	keyEx := expression.Key("golfer_id").Equal(expression.Value(golferID))
	filterEx := expression.Name("played_at").Between(expression.Value(startDate), expression.Value(endDate))

	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).WithFilter(filterEx).Build()
	if err != nil {
		return scores, nil, err
	}

	// Assemble Query
	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(c.ScoresTableName),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		IndexName:                 aws.String("golfer_id-index"),
	}

	if pageKey != nil {
		queryInput.ExclusiveStartKey = pageKey
	}

	result, err := c.dynamoClient.Query(c.Context, queryInput)
	if err != nil {
		return scores, nil, err
	}

	// Deserialize each item in the result
	for _, item := range result.Items {
		var score ghin.Score
		err = attributevalue.UnmarshalMap(item, &score)
		if err != nil {
			return scores, nil, err
		}

		scores = append(scores, score)
	}

	if result.LastEvaluatedKey != nil {
		// Results are paginated
		pageKey = result.LastEvaluatedKey
	}

	return scores, pageKey, nil
}

// UpdateScore updates a Score
func (c *Client) UpdateScore(score ghin.Score) (ghin.Score, error) {
	update, err := attributevalue.MarshalMap(score)
	if err != nil {
		return score, err
	}

	_, err = c.dynamoClient.PutItem(c.Context, &dynamodb.PutItemInput{
		TableName: aws.String(c.ScoresTableName),
		Item:      update,
	})

	if err != nil {
		log.Printf("failed to update Score: %+v\nScore: %s", err, hcutil.ObjectToJSON(score))
	}

	return score, err
}

// DeleteScore deletes a Score given its ID
func (c *Client) DeleteScore(scoreID string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(c.ScoresTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: scoreID,
			},
		},
	}

	_, err := c.dynamoClient.DeleteItem(c.Context, input)

	if err != nil {
		log.Printf("failed to delete Score for Score (%s): %+v", scoreID, err)
	}

	return err
}
