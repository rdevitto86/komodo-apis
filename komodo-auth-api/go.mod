module komodo-auth-api

go 1.25

require github.com/go-chi/chi/v5 v5.2.3

require (
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/google/uuid v1.6.0
	komodo-internal-lib-apis-go v0.0.0-00010101000000-000000000000
// komodo-internal-lib-apis-go v0.1.0
)

replace komodo-internal-lib-apis-go => ../komodo-internal-lib-apis-go

require (
	github.com/aws/aws-sdk-go-v2 v1.39.2 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.31.11 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.18.15 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.39.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.29.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.38.6 // indirect
	github.com/aws/smithy-go v1.23.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/redis/go-redis/v9 v9.15.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
