# Komodo E-Commerce Monorepo

Production-style e-commerce platform. Independently deployable Go microservices, a SvelteKit SSR frontend, shared SDKs, and local AWS infrastructure via LocalStack.

---

## Getting Started

**Prerequisites:** Docker >= 24.x, Go 1.26+, Bun (for SvelteKit/TS SDK), Make

```bash
# Start local backing services (DynamoDB, S3, Secrets Manager, Redis)
cd localstack && docker compose up -d

# Run any service
cd komodo-<service> && docker compose up --build
```

---

## Conventions

| Concern | Convention |
|---------|-----------|
| API style | JSON over HTTP, REST |
| Routing | Go 1.26 `net/http` ServeMux (`GET /path/{id}`) |
| Auth | JWT (RS256) via `komodo-auth-api`. Service-to-service via client credentials. |
| Logging | `slog` structured JSON. `tint` locally, JSON in staging/prod. |
| Errors | RFC 7807 Problem+JSON |
| Schema | `docs/openapi.yaml` per service |
| Tracing | OpenTelemetry OTLP (planned) |
| Secrets | AWS Secrets Manager via `komodo-forge-sdk-go` at startup |
| Networking | All services share `komodo-network` (created by LocalStack compose) |

---

## Services

| Port | Service | Domain | Status |
|------|---------|--------|--------|
| 7001 | `komodo-ssr-engine-svelte` | Frontend & Infra | Active |
| 7002 | `komodo-core-cdn-api` | Frontend & Infra | Stub |
| 7011 | `komodo-auth-api` | Identity & Security | Active |
| 7021 | `komodo-core-entitlements-api` | Core Platform | Stub |
| 7022 | `komodo-core-features-api` | Core Platform | Stub |
| 7031 | `komodo-address-api` | Address & Geo | Active |
| 7041 | `komodo-shop-items-api` | Commerce & Catalog | Active |
| 7042 | `komodo-search-api` | Commerce & Catalog | Stub |
| 7051 | `komodo-user-api` | User & Profile | Active |
| 7061 | `komodo-order-api` | Orders | Scaffolded |
| 7071 | `komodo-payments-api` | Payments | Scaffolded |
| 7081 | `komodo-communications-api` | Communications | Scaffolded |
| 7091 | `komodo-loyalty-api` | Loyalty & Social | Scaffolded |
| 7092 | `komodo-reviews-api` | Loyalty & Social | Scaffolded |
| 7101 | `komodo-support-api` | Support & CX | Scaffolded |
| 7111 | `komodo-analytics-collector-api` | Analytics | Stub |

> Port override: set the `PORT` env var on any service.

### Service Status
- **Active** — Implemented and running
- **Scaffolded** — Directory structure exists, not yet implemented
- **Stub** — Empty module or `main.go` only

---

## Shared Libraries

**`komodo-forge-sdk-go`** — Internal Go SDK. AWS clients (DynamoDB, S3, Secrets Manager, Redis, Aurora), HTTP middleware stack, JWT/JWKS crypto, structured logging, concurrency utilities.

**`komodo-forge-sdk-ts`** — Internal TypeScript SDK. Domain types, API client utilities, frontend helpers. Backend modules (logging, telemetry) are stubs.

---

## Infrastructure

LocalStack (`localstack/`) emulates AWS locally:

| Service | Purpose |
|---------|---------|
| Secrets Manager | Service secrets (DB passwords, API keys, JWT keys) |
| S3 | Product data, content, file storage |
| DynamoDB | User data (NoSQL) |
| RDS | Aurora-compatible relational DB (planned) |
| Redis | Sessions and caching (standalone container, port 6379) |

Init scripts in `localstack/init/` seed all services on startup.
