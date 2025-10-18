package secretsmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"komodo-internal-lib-apis-go/config"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var (
	secretsManagerClient *secretsmanager.Client
	useAWS               bool
	secretPrefix         string // e.g., "some/prod/"
	batchSecretName      string // Name of the batch secret containing all keys as JSON
)

func init() {
	// Default secret prefix
	if secretPrefix = config.GetConfigValue("AWS_SECRET_PREFIX"); secretPrefix == "" {
		logger.Error("failed to set AWS Secrets prefix")
		return
	}

	// Default batch secret name
	if batchSecretName = config.GetConfigValue("AWS_BATCH_SECRET_NAME"); batchSecretName == "" {
		batchSecretName = "all-secrets"
	}

	// Check if AWS Secrets Manager should be used
	if useAWS = config.GetConfigValue("USE_AWS_SECRETS") == "true"; useAWS {
		// Initialize AWS Secrets Manager client
		cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
			awsconfig.WithRegion(config.GetConfigValue("AWS_REGION")),
		)
		if err != nil {
			logger.Error("failed to load AWS config", err)
			useAWS = false
			return
		}

		secretsManagerClient = secretsmanager.NewFromConfig(cfg)
		logger.Info("AWS Secrets Manager client initialized")
	} else {
		logger.Info("AWS Secrets Manager failed to initialize or is disabled")
	}
}

// Loads multiple secrets using AWS Secrets Manager or from .env.secrets
func LoadSecrets(keys []string) error {
	if !IsUsingAWS() {
		logger.Error("AWS Secrets Manager not configured")
		return fmt.Errorf("AWS Secrets Manager not configured")
	}
	_, err := GetSecrets(keys)
	return err
}

// GetSecret retrieves a single secret from AWS Secrets Manager or .env.secrets
func GetSecret(key string) (string, error) {
	if IsUsingAWS() {
		secretName := secretPrefix + key // e.g., "some-api/prod/JWT_PRIVATE_KEY"
		input := &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretName),
		}

		result, err := secretsManagerClient.GetSecretValue(context.TODO(), input)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to retrieve secret %s from AWS", secretName), err)
			return "", err
		}
		if result.SecretString == nil {
			return "", fmt.Errorf("secret %s has no string value", secretName)
		}

		logger.Info(fmt.Sprintf("successfully retrieved secret %s from AWS", key))
		config.SetConfigValue(key, *result.SecretString)
		return *result.SecretString, nil
	}
	return "", fmt.Errorf("AWS Secrets Manager not configured")
}

// Retrieves multiple secrets using AWS batch call (single JSON secret) or from local secrets
// Returns a map of key-value pairs
func GetSecrets(keys []string) (map[string]string, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("no keys provided for secret retrieval")
	}
	if IsUsingAWS() {
		secretName := secretPrefix + batchSecretName // e.g., "some-api/prod/all-secrets"
		input := &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretName),
		}

		result, err := secretsManagerClient.GetSecretValue(context.TODO(), input)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to retrieve batch secret %s from AWS", secretName), err)
			return nil, err
		}
		if result.SecretString == nil {
			return nil, fmt.Errorf("batch secret %s has no string value", secretName)
		}

		// Parse JSON
		var allSecrets map[string]string
		if err := json.Unmarshal([]byte(*result.SecretString), &allSecrets); err != nil {
			logger.Error("failed to parse batch secret JSON", err)
			return nil, err
		}

		// Filter only requested keys
		secrets := make(map[string]string)
		missingKeys := []string{}

		for _, key := range keys {
			if value, ok := allSecrets[key]; ok {
				secrets[key] = value
				config.SetConfigValue(key, value)
			}
		}

		if len(missingKeys) > 0 {
			logger.Warn(fmt.Sprintf("keys not found in batch secret: %v", missingKeys))
		}

		logger.Info(fmt.Sprintf("successfully retrieved %d secrets from AWS batch", len(secrets)))
		return secrets, nil
	}
	return nil, fmt.Errorf("AWS Secrets Manager not configured")
}

// Checks if the client is configured to use AWS Secrets Manager
func IsUsingAWS() bool { return useAWS && secretsManagerClient != nil }
