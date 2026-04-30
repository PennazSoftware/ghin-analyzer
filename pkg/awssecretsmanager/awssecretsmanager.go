package awssecretsmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// AwsSecretsManager interface describes the methods necessary for implementing Secrets Manager client
type AwsSecretsManager interface {
	GetSecret(secretName string, secretObjectPtr interface{}) error
}

// verifying if the Client struct is indeed implementing the AwsSecretsManager interface
var _ AwsSecretsManager = (*DefaultAwsSecretsManager)(nil)

// DefaultAwsSecretsManager is the struct that embodies the AWS Secrets Manager client
type DefaultAwsSecretsManager struct {
	secretsManager *secretsmanager.Client
	context        context.Context
}

// GHINSecret contains the secret info for accessing GHIN APIs
type GHINSecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// New creates a new AWS Secrets Manager client
func New() *DefaultAwsSecretsManager {
	ctx := context.TODO()

	// Load the SDK configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"),
	)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create the Secrets Manager client
	secretsClient := secretsmanager.NewFromConfig(cfg)

	return &DefaultAwsSecretsManager{
		secretsManager: secretsClient,
		context:        ctx,
	}
}

// GetSecret retrieves the secret given the secret name and the interface type the secret
// should be cast to
func (c *DefaultAwsSecretsManager) GetSecret(secretName string, secretObjectPtr interface{}) error {
	secretValueInput := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	// Retrieve the secret
	result, err := c.secretsManager.GetSecretValue(c.context, secretValueInput)
	if err != nil {
		return fmt.Errorf("failed to retrieve secret with name of '%s'. %s", secretName, err.Error())
	}

	// unmarshal to the required type
	if result.SecretString != nil {
		err = json.Unmarshal([]byte(*result.SecretString), &secretObjectPtr)
		if err != nil {
			return fmt.Errorf("error unmarshalling the secret (%s) to desired type. %s", secretName, err.Error())
		}
	} else {
		b, ok := secretObjectPtr.(*[]byte)
		if ok {
			*b = append(*b, result.SecretBinary...)
		} else {
			return fmt.Errorf("error casting interface to byte for binary secret")
		}
	}

	return nil
}
