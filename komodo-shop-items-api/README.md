# komodo-shop-items-api

Consolidated shop items API for the Komodo e-commerce platform. Serves product/service catalog data and inventory from S3, with authenticated suggestions powered by a recommendation engine.

## Routes

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | None | Health check |
| GET | `/item/inventory` | Public | Bulk inventory/stock status |
| GET | `/item/{sku}` | Public | Single product or service by SKU |
| POST | `/item/suggestion` | Bearer JWT | Personalized product suggestions |

## Build

```bash
make build                # Build Docker image (local)
make build ENV=dev        # Build for dev environment
make build ENV=prod       # Build for production
```

## Run

```bash
make bootstrap            # Build + run (default: local)
make bootstrap ENV=dev    # Dev environment
make run                  # Run only (image must exist)
make stop                 # Stop and remove containers
make restart              # Stop + run
```

## Misc

| Key | Value |
|-----|-------|
| **Language** | Go 1.26 |
| **Router** | go-chi/chi v5 |
| **Port** | 7050 |
| **Data Store** | AWS S3 (product/service JSON, inventory manifest) |
| **Secrets** | AWS Secrets Manager (LocalStack for local dev) |
| **SDK** | komodo-forge-sdk-go (S3 client, shop_items models/adapters) |
| **Health** | `GET /health` → `{"status": "OK"}` |
| **Docs** | `docs/openapi.yaml` |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `APP_NAME` | Service name (`komodo-shop-items-api`) |
| `ENV` | Environment (`local`, `dev`, `staging`, `prod`) |
| `PORT` | HTTP listen port (default `7050`) |
| `AWS_REGION` | AWS region |
| `AWS_ENDPOINT` | AWS/LocalStack endpoint |
| `AWS_SECRET_PREFIX` | Secrets Manager key prefix |
| `AWS_SECRET_BATCH` | Secrets Manager batch secret name |
| `S3_ENDPOINT` | S3 endpoint (from Secrets Manager) |
| `S3_ACCESS_KEY` | S3 access key (from Secrets Manager) |
| `S3_SECRET_KEY` | S3 secret key (from Secrets Manager) |
| `S3_ITEMS_BUCKET` | S3 bucket for product/service/inventory JSON |

## S3 Bucket Layout

```
s3://<S3_ITEMS_BUCKET>/
├── products/<sku>.json       # Product JSON per SKU
├── services/<sku>.json       # Service JSON per SKU
└── inventory/manifest.json   # Inventory manifest (all tracked items)
```

## Testing

```bash
make test           # All tests with race detector
make test_unit      # Unit tests only
make test_e2e       # End-to-end tests
```
