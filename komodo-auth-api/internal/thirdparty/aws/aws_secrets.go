package aws

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var SecretsManagerClient *secretsmanager.Client

func InitSecretsClient() {
	env := os.Getenv("API_ENV")
	if env != "prod" && env != "staging" {
		log.Println("InitSecretsClient: skipping Elasticache initialization in local/DEV environment")
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO()) // TODO configure AWS
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	SecretsManagerClient = secretsmanager.NewFromConfig(cfg)
}

func GetSecret(secretName string) (string, error) {
	if SecretsManagerClient == nil {
		return "", errors.New("SecretsManager client is not initialized")
	}

	res, err := SecretsManagerClient.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})

	if err != nil { return "", err }
	return *res.SecretString, nil
}