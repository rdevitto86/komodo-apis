package aws

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var SecretsManagerClient *secretsmanager.Client

func InitSecretsClient() {
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