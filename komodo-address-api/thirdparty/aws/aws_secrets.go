package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetSecret(secretName string) (string, error) {
	ctx := context.Background()
	awsCfg, err := config.LoadDefaultConfig(ctx)

	if err != nil { return "", err }

	client := secretsmanager.NewFromConfig(awsCfg)
	out, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	})

	if err != nil { return "", err }
	return *out.SecretString, nil
}
