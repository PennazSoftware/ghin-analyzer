package main

import (
	"PennazSoftware/ghin-analyzer/pkg/awsdynamodb"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGolfer(t *testing.T) {
	os.Setenv("AWS_PROFILE", "pennaz")
	// Initialize the AWS DynamoDB client and GHIN API client
	dbClient := awsdynamodb.New("dev")

	golfer, err := dbClient.GetGolfer(50155)
	assert.NoError(t, err)
	assert.Equal(t, 50155, golfer.GolferID)
}
