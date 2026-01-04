package secretsmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"komodo-forge-apis-go/config"
	logger "komodo-forge-apis-go/logging/runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var secretsManagerClient *secretsmanager.Client

type Config struct {
	Region		string
	Endpoint 	string
	Prefix   	string
	Batch 		string
	Keys			[]string
}

// Initialize AWS Secrets Manager and load secrets in one call
func Bootstrap(cfg Config) error {
	if cfg.Region == "" {
		logger.Error("region not provided")
		return fmt.Errorf("aws region not provided for secrets manager")
	}

	var awsCfg aws.Config
	var err error
 
	// Load AWS config
	cfgOpts := []func(*awsconfig.LoadOptions) error{awsconfig.WithRegion(cfg.Region)}
	
	// For LocalStack, provide dummy credentials to avoid EC2 IMDS lookup
	if cfg.Endpoint != "" {
		cfgOpts = append(cfgOpts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("test", "test", ""),
		))
	}
	
	awsCfg, err = awsconfig.LoadDefaultConfig(context.TODO(), cfgOpts...)
	if err != nil {
		logger.Error("failed to load AWS config", err)
		return err
	}

	// Override endpoint if provided (for LocalStack)
	if cfg.Endpoint != "" {
		secretsManagerClient = secretsmanager.NewFromConfig(awsCfg, func(opts *secretsmanager.Options) {
			opts.BaseEndpoint = aws.String(cfg.Endpoint)
		})
		logger.Info("aws secrets manager client initialized with custom endpoint: " + cfg.Endpoint)
	} else {
		secretsManagerClient = secretsmanager.NewFromConfig(awsCfg)
		logger.Info("aws secrets manager client initialized")
	}

	if len(cfg.Keys) > 0 {
		if _, err := GetSecrets(cfg.Keys, cfg.Prefix, cfg.Batch); err != nil {
			logger.Error("failed to load secrets", err)
			return err
		}
	}
	return nil
}

// Retrieves a single secret
func GetSecret(key string, prefix string) (string, error) {
	if secretsManagerClient == nil {
		logger.Error("aws secrets manager client not initialized")
		return "", fmt.Errorf("aws secrets manager client not initialized")
	}
	if prefix == "" {
		logger.Error("secret prefix not initialized")
		return "", fmt.Errorf("secret prefix not initialized")
	}

	secretPath := prefix + key // e.g., "some-api/prod/JWT_PRIVATE_KEY"
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretPath),
	}

	result, err := secretsManagerClient.GetSecretValue(context.TODO(), input)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to retrieve secret %s from AWS", secretPath), err)
		return "", err
	}
	if result.SecretString == nil {
		return "", fmt.Errorf("secret %s has no string value", secretPath)
	}

	logger.Info(fmt.Sprintf("successfully retrieved secret %s from AWS", key))
	config.SetConfigValue(key, *result.SecretString)
	return *result.SecretString, nil
}

// Retrieves multiple secrets using AWS batch call
func GetSecrets(keys []string, prefix string, batchId	string) (map[string]string, error) {
	if secretsManagerClient == nil {
		logger.Error("aws secrets manager client not initialized")
		return nil, fmt.Errorf("aws secrets manager client not initialized")
	}
	if len(keys) == 0 {
		logger.Warn("no secrets to load")
		return nil, fmt.Errorf("no keys provided for secret retrieval")
	}
	if prefix == "" {
		logger.Error("secret prefix not initialized")
		return nil, fmt.Errorf("secret prefix not initialized")
	}
	if batchId == "" {
		logger.Error("batch secret name not initialized")
		return nil, fmt.Errorf("batch secret name not initialized")
	}

	secretPath := prefix + batchId // e.g., "some-api/prod/all-secrets"
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretPath),
	}

	result, err := secretsManagerClient.GetSecretValue(context.TODO(), input)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to retrieve batch secret %s from AWS", secretPath), err)
		return nil, err
	}
	if result.SecretString == nil {
		logger.Error(fmt.Sprintf("batch secret %s has no string value", secretPath))
		return nil, fmt.Errorf("batch secret %s has no string value", secretPath)
	}

	// Parse JSON string
	var allSecrets map[string]string
	if err := json.Unmarshal([]byte(*result.SecretString), &allSecrets); err != nil {
		logger.Error("failed to parse batch secret response for " + secretPath, err)
		return nil, err
	}

	secrets := make(map[string]string)
	missingKeys := []string{}

	for _, key := range keys {
		if value, ok := allSecrets[key]; ok {
			secrets[key] = value
			config.SetConfigValue(key, value)
		} else {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		logger.Warn(fmt.Sprintf("keys not found in batch secret: %v", missingKeys))
	}

	logger.Info(fmt.Sprintf("successfully retrieved %d secrets from AWS batch", len(secrets)))
	return secrets, nil
}
