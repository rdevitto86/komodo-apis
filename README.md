# Komodo E-Commerce Monorepo

A monorepo for a production-style e-commerce platform containing independently deployable APIs, a SvelteKit SSR frontend, shared SDKs, and local infrastructure.

---

## Table of Contents

- [Getting Started](#getting-started)
- [Local Development](#local-development)
- [Conventions](#conventions)
- [Port Allocation](#port-allocation)
- [Services](#services)
  - [komodo-auth-api](#komodo-auth-api)
  - [komodo-address-api](#komodo-address-api)
  - [komodo-shop-items-api](#komodo-shop-items-api)
  - [komodo-user-api](#komodo-user-api)
  - [komodo-ssr-engine-svelte](#komodo-ssr-engine-svelte)
- [Shared Libraries](#shared-libraries)
  - [komodo-forge-sdk-go](#komodo-forge-sdk-go)
  - [komodo-forge-sdk-ts](#komodo-forge-sdk-ts)
- [Planned Services](#planned-services)
- [Infrastructure](#infrastructure)

---

## Getting Started

**Prerequisites**
- Docker >= 24.x, Docker Compose
- Go 1.26+
- Node.js LTS / Bun (for SvelteKit and TS SDK)
- Make (optional but recommended)

**Bootstrap**
```bash
git clone <your-fork-url> komodo-apis && cd komodo-apis
```

**Run LocalStack (backing services)**
```bash
cd localstack && docker compose up -d
```

**Run a service**
```bash
cd komodo-auth-api && docker compose up --build
```

---

## Local Development

- **Ports:** Each service exposes a unique port (see [Port Allocation](#port-allocation)). Override via `PORT` env var.
- **Backing Services:** LocalStack provides Secrets Manager, S3, DynamoDB, and RDS locally. Redis runs as a standalone container.
- **Secrets:** All implemented services pull secrets from AWS Secrets Manager at startup via `komodo-forge-sdk-go`.
- **Networking:** Services share a Docker network (`komodo-network`) created by the LocalStack compose file.

---

## Conventions

- **API Style:** JSON over HTTP, REST.
- **Auth:** JWT access tokens via `komodo-auth-api`. Service-to-service auth via client credentials.
- **Tracing:** OpenTelemetry (OTLP) — planned, stubs in SDK.
- **Logging:** Structured JSON via `slog` (Go) to stdout. Pretty-printed locally via `tint`, JSON in staging/prod.
- **Errors:** RFC 7807 (Problem+JSON) recommended.
- **Schema:** OpenAPI spec per service in `docs/openapi.yaml`.

---

## Port Allocation

Services are allocated port ranges by business domain.

> **Note:** Port 7000 is reserved due to macOS system service conflicts. All services start at 7001.

| Port Range | Domain | Priority | Description |
|------------|--------|----------|-------------|
| 7001-7019 | Infrastructure & Mesh | CRITICAL | Sidecars, Service Discovery, Config management |
| 7020-7049 | Identity & Security | CRITICAL | Auth, IAM, Token validation |
| 7050-7149 | Core Business Logic | HIGH | Inventory, Orders, Payments |
| 7150-7249 | User & Engagement | MEDIUM | Profiles, Loyalty, Communications |
| 7250-7349 | Emerging Tech & AI | MEDIUM | GenAI, LLM Orchestration |
| 7350-7449 | Observability & Ops | LOW | Analytics collector, health dashboards |

### Assigned Ports

| Port | Service | Status |
|------|---------|--------|
| 7010 | `komodo-address-api` | Implemented |
| 7010 | `komodo-ssr-engine-svelte` | Implemented |
| 7020 | `komodo-auth-api` | Implemented |
| 7050 | `komodo-shop-items-api` | Implemented |
| 7150 | `komodo-user-api` | Implemented |

---

## Services

### komodo-auth-api

OAuth 2.0 authorization server. Issues and validates JWT access tokens, exposes a JWKS endpoint, and supports token introspection and revocation.

| | |
|---|---|
| **Language** | Go (Chi) |
| **Port** | 7020 |
| **Data Stores** | ElastiCache (Redis), Secrets Manager |
| **Docker** | `docker compose up --build` |
| **Health** | `GET /health` |

**Routes**
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | No | Health check |
| GET | `/.well-known/jwks.json` | No | Public JWKS endpoint |
| POST | `/oauth/token` | No | Issue access token |
| GET | `/oauth/authorize` | No | Authorization flow |
| POST | `/oauth/introspect` | Yes | Token introspection |
| POST | `/oauth/revoke` | Yes | Token revocation |

**Middleware:** Request ID, Telemetry, Rate Limiter, IP Access, Security Headers, Normalization, Sanitization, Rule Validation, Client Type, Auth (protected routes).

---

### komodo-address-api

Address validation, normalization, and geocoding service.

| | |
|---|---|
| **Language** | Go (Gin) |
| **Port** | 7010 |
| **Data Stores** | Secrets Manager |
| **Health** | `GET /health` |

**Routes**
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | Yes | Health check |
| POST | `/validate` | Yes | Validate an address |
| POST | `/normalize` | Yes | Normalize an address |
| POST | `/geocode` | Yes | Geocode an address |

> **Note:** This service uses Gin and a different auth pattern (external token validation URL) compared to the other Go services which use Chi + the forge SDK middleware stack. Consider migrating to the shared SDK pattern.

---

### komodo-shop-items-api

Product catalog and inventory service. Serves item data from S3 and supports authenticated product suggestions.

| | |
|---|---|
| **Language** | Go (Chi) |
| **Port** | 7050 |
| **Data Stores** | S3, Secrets Manager |
| **Docker** | `docker compose up --build` |
| **Health** | `GET /health` |

**Routes**
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | No | Health check |
| GET | `/item/inventory` | No | Get inventory listing |
| GET | `/item/{sku}` | No | Get item by SKU |
| POST | `/item/suggestion` | Yes | Get product suggestions |

**Middleware:** Request ID, Telemetry, Rate Limiter, IP Access, CORS, Security Headers. Protected routes add Auth, CSRF, Normalization, Sanitization, Rule Validation.

---

### komodo-user-api

Core user service. Manages profiles, addresses, orders, payments, and preferences for authenticated users.

| | |
|---|---|
| **Language** | Go (Chi) |
| **Port** | 7150 |
| **Data Stores** | DynamoDB, Secrets Manager |
| **Docker** | `docker compose up --build` |
| **Health** | `GET /health` |

**Routes**
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | No | Health check |
| POST | `/me/profile` | Yes | Get profile |
| PUT | `/me/profile` | Yes | Update profile |
| DELETE | `/me/profile` | Yes | Delete profile |
| POST | `/me/profile/create` | Yes | Create user |
| POST | `/me/addresses/query` | Yes | Get addresses |
| POST | `/me/addresses/create` | Yes | Add address |
| PUT | `/me/addresses/update` | Yes | Update address |
| DELETE | `/me/addresses/delete` | Yes | Delete address |
| POST | `/me/orders` | Yes | Get orders |
| PUT | `/me/orders` | Yes | Update order |
| POST | `/me/orders/create` | Yes | Create order |
| POST | `/me/orders/cancel` | Yes | Cancel order |
| POST | `/me/orders/return` | Yes | Return order |
| POST | `/me/payments` | Yes | Get payments |
| PUT | `/me/payments` | Yes | Upsert payment |
| DELETE | `/me/payments` | Yes | Delete payment |
| GET | `/me/preferences` | Yes | Get preferences |
| PUT | `/me/preferences` | Yes | Update preferences |

**Middleware:** Request ID, Telemetry, Rate Limiter, IP Access, CORS, Security Headers, Auth, CSRF, Normalization, Rule Validation, Sanitization, Idempotency.

---

### komodo-ssr-engine-svelte

Server-side rendered SvelteKit frontend with a backend-for-frontend (BFF) API layer. Serves pages and proxies API calls to downstream services.

| | |
|---|---|
| **Language** | TypeScript (SvelteKit 5, Vite) |
| **Port** | 7010 |
| **Data Stores** | S3, CloudFront |
| **Runtime** | Node (adapter-node) |
| **Docker** | `docker compose up --build` |

**Pages:** `/landing`, `/products`, `/products/[id]`, `/orders`, `/orders/[id]`, `/services`, `/services/[id]`, `/marketing`, `/about`, `/contact`, `/faq`, `/terms`

**BFF API Routes** (under `/api/v1`):
- `/health` — Health check
- `/landing` — Landing page data
- `/products`, `/products/[id]` — Product data
- `/orders`, `/orders/[id]` — Order data
- `/services`, `/services/[id]` — Service data
- `/services/scheduling`, `/services/scheduling/[id]` — Scheduling
- `/marketing/content`, `/marketing/content/[id]` — Marketing content
- `/marketing/user`, `/marketing/user/[id]` — User marketing
- `/admin/manage/content/upsert` — Admin content management
- `/admin/manage/content/invalidate` — Admin cache invalidation

**Dependencies:** `@komodo-forge-sdk/typescript` (linked), `@aws-sdk/client-s3`, `@aws-sdk/client-cloudfront`

---

## Shared Libraries

### komodo-forge-sdk-go

Internal Go SDK shared across all Go microservices. Provides:

- **AWS clients:** DynamoDB, S3, Secrets Manager, ElastiCache (Redis), Aurora
- **HTTP middleware:** Auth, CORS, CSRF, Rate Limiter, IP Access, Idempotency, Normalization, Sanitization, Rule Validation, Security Headers, Request ID, Telemetry, Redaction, Client Type
- **Crypto:** JWT (RS256 sign/verify, JWKS), OAuth client
- **Logging:** Structured `slog` logger with env-aware formatting (tint for local, JSON for deployed) and field redaction
- **Config:** Centralized config value resolution
- **Concurrency:** Semaphore and worker pool utilities

### komodo-forge-sdk-ts

Internal TypeScript SDK shared by the SvelteKit frontend and any Node-based services. Provides:

- **Backend:** AWS clients, config, DB utilities, logging (runtime, security, telemetry), middleware, observability
- **Frontend:** API client utilities
- **Shared:** Crypto, domain types (auth, payments, user, orders, products, services, marketing), entitlements, feature flags, security, utilities

> **Note:** Most backend modules (logging, telemetry, observability) are currently stubs awaiting implementation.

---

## Planned Services

The following services exist as empty scaffolds (empty `cmd/` and `internal/` directories) or stubs (empty `main.go` files). They are not yet implemented.

### Scaffolded (empty directory structure)

| Service | Purpose |
|---------|---------|
| `komodo-order-api` | Order management |
| `komodo-payments-api` | Payment processing |
| `komodo-communications-api` | Email, SMS, push notifications |
| `komodo-loyalty-api` | Loyalty programs and rewards |
| `komodo-reviews-api` | Product and service reviews |
| `komodo-support-api` | Customer support / ticketing |

### Stubs (empty Go modules or package.json)

| Service | Purpose |
|---------|---------|
| `komodo-analytics-collector-api` | Clickstream and business event ingestion (Lambda target) |
| `komodo-core-entitlements-api` | User/account entitlements |
| `komodo-core-features-api` | Feature flag management |
| `komodo-core-files-api` | File upload/management (TypeScript) |
| `komodo-search-api` | Product/content search |

---

## Infrastructure

### LocalStack (`localstack/`)

Docker Compose setup providing local AWS service emulation:

- **Secrets Manager** — Stores service secrets (DB passwords, API keys, JWT keys)
- **S3** — Object storage for product data, content, files
- **DynamoDB** — NoSQL store for user data
- **RDS** — Relational DB (Aurora compatible)
- **Redis** — Standalone container (port 6379) for caching and session storage

**Start:**
```bash
cd localstack && docker compose up -d
```

Init scripts in `localstack/init/` automatically seed Secrets Manager, S3, and DynamoDB on startup.
