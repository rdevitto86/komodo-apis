package secretsmanager

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var secretsManagerClient *secretsmanager.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO()) // TODO configure AWS
	if err != nil {
		log.Fatalf("[FATAL] failed to load config: %v", err)
	}

	secretsManagerClient = secretsmanager.NewFromConfig(cfg)
}

func GetSecrets(keys []string) (map[string]string, error) {
	if secretsManagerClient == nil {
		return nil, errors.New("[ERROR] SecretsManager client is not initialized")
	}

	// TODO batch call

	return nil, nil
}

func SetSecrets(keys []string) {
	if secretsManagerClient == nil {
		log.Println("[WARN] SetSecrets: SecretsManager client is not initialized")
		return
	}

	// TODO batch call
}