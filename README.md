# Komodo E-Commerce Backend Monorepo

A single repo showcasing a complete, production-style backend for a modern e-commerce platform. It contains independently deployable APIs/services, shared libs, and infra that you can build/run locally or in containers. This README is a template you can keep up-to-date as you implement each service.

> Tip: Every API section below includes **Build**, **Run**, and **Misc** (language/stack, operations, etc.). Replace any `👈 fill me in` placeholders as you wire things up.

---

## Table of Contents

- [Monorepo Structure](#monorepo-structure)  
- [Getting Started (All Services)](#getting-started-all-services)  
- [Local Development](#local-development)  
- [Conventions](#conventions)  
- [Services](#services)
  - [komodo-address-api](#komodo-address-api)
  - [komodo-ai-chatbot-api](#komodo-ai-chatbot-api)
  - [komodo-ai-summary-api](#komodo-ai-summary-api)
  - [komodo-analytics-interaction-api](#komodo-analytics-interaction-api)
  - [komodo-analytics-logs-api](#komodo-analytics-logs-api)
  - [komodo-analytics-telemetry-api](#komodo-analytics-telemetry-api)
  - [komodo-auth-api](#komodo-auth-api)
  - [komodo-catalog-api](#komodo-catalog-api)
  - [komodo-entitlements-api](#komodo-entitlements-api)
  - [komodo-knowledge-api](#komodo-knowledge-api)
  - [komodo-order-api](#komodo-order-api)
  - [komodo-order-discount-api](#komodo-order-discount-api)
  - [komodo-order-scheduling-api](#komodo-order-scheduling-api)
  - [komodo-payments-api](#komodo-payments-api)
  - [komodo-search-api](#komodo-search-api)
  - [komodo-servicing-chat-api](#komodo-servicing-chat-api)
  - [komodo-ssr-engine](#komodo-ssr-engine)
  - [komodo-user-api](#komodo-user-api)
  - [komodo-user-marketing-api](#komodo-user-marketing-api)
  - [komodo-user-recommendations-api](#komodo-user-recommendations-api)
  - [komodo-user-reviews-api](#komodo-user-reviews-api)
- [Testing, Observability & Security](#testing-observability--security)
- [CI/CD](#cicd)
- [Troubleshooting](#troubleshooting)
- [License](#license)

---

## Getting Started (All Services)

**Prereqs**
- Docker ≥ 24.x, Docker Compose  
- Node.js LTS (if using Node services) **or** language runtimes per service  
- Make (optional but recommended)  
- `just` or `make` (choose one; commands below show `make`)  

**Bootstrap**
```bash
git clone <your-fork-url> komodo && cd komodo
cp .env.example .env
make bootstrap   # installs toolchains, hooks, and package deps across the monorepo
```

**Run Everything (happy path)**
```bash
make up          # docker compose up all core services + backing stores
make logs        # tail logs
make down        # stop
```

**Run One Service Only**
```bash
make up SERVICE=komodo-catalog-api
```

---

## Local Development

- **Ports:** Each service exposes a unique port. Defaults are listed in each section below (override via env).  
- **Configuration:** Put shared config in `.env`. Service-specific variables live in `services/<name>/.env`.  
- **Databases:** Postgres and Redis containers are provided in `docker-compose.yml`.  
- **Migrations:** Standardized via `make migrate SERVICE=<name>`.  
- **Seeds:** `make seed SERVICE=<name>`.

---

## Conventions

- **API Style:** JSON over HTTP; REST first; gRPC optional for internal calls.  
- **Auth:** JWT access tokens via `komodo-auth-api`. Service-to-service via mTLS or signed service tokens.  
- **Tracing:** OpenTelemetry (OTLP) → collector (Jaeger/Tempo).  
- **Logging:** Structured JSON, request IDs, 12-factor log to stdout.  
- **Errors:** RFC 7807 (Problem+JSON) recommended.  
- **Schema:** OpenAPI per service at `/docs` and `/openapi.json`.

---

## Services

Each service uses the same template:

- **Overview**  
- **Build**  
- **Run**  
- **Misc** (Language/Stack, Port, Data Stores, Key Operations, Docs, Health)  

Replace `👈 fill me in` as you implement.

---

### komodo-address-api

**Overview**  
Validates, normalizes, and geocodes shipping/billing addresses.

**Misc**  
- **Language/Stack:** Golang
- **Port:** 7001  
- **Data Stores:** Postgres, Redis
- **Key Ops:** Validate, Normalize, Geocode  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-ai-chatbot-api

**Overview**  
Provides LLM generated servicing via self-service chat

**Misc**  
- **Language/Stack:** Python + FastAPI
- **Port:** 7002  
- **Upstreams:** LLM provider(s)
- **Key Ops:** Message
- **Docs/Health:** `/docs`, `/health`

---

### komodo-ai-summary-api

**Overview**  
Generates product/order/user summaries (LLMs).

**Misc**  
- **Language/Stack:** Python + FastAPI
- **Port:** 7022  
- **Upstreams:** LLM provider(s)  
- **Key Ops:** Summaries, Models listing  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-analytics-interaction-api

**Overview**  
Captures user interactions for real-time personalization.

**Misc**  
- **Language/Stack:** Golang
- **Port:** 7003  
- **Data:** Kafka/Redpanda  
- **Key Ops:** Events ingestion, Stats query  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-analytics-logs-api

**Overview**  
Ingests application logs and exposes query endpoints.

**Misc**  
- **Language/Stack:** Golang
- **Port:** 7004  
- **Data:** Loki/Elastic  
- **Key Ops:** Ingest, Query logs  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-analytics-telemetry-api

**Overview**  
Metrics/trace ingestion + query proxy.

**Misc**  
- **Language/Stack:** Golang
- **Port:** 7005  
- **Data:** Prometheus, Jaeger/Tempo  
- **Key Ops:** Metrics, Traces, Dashboards  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-auth-api

**Overview**  
Authentication, authorization, MFA, etc.

**Misc**  
- **Language/Stack:** Golang
- **Port:** 7017  
- **Key Ops:** Signup, Login, Token refresh, Password reset  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-catalog-api

**Overview**  
Products, variants, categories, pricing, inventory.

**Misc**  
- **Language/Stack:** Golang
- **Port:** 7006  
- **Data:** Postgres, Redis, S3  
- **Key Ops:** Products CRUD, Inventory, Categories  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-entitlements-api

**Overview**  
User entitlements and feature toggles

**Misc**  
- **Language/Stack:** Golang
- **Port:** 7007  
- **Key Ops:** Entitlements, Toggles
- **Docs/Health:** `/docs`, `/health`

---

### komodo-knowledge-api

**Overview**  
Knowledge base and FAQ retrieval.

**Misc**  
- **Language/Stack:** Node + Fastify + Prisma
- **Port:** 7008  
- **Data:** Postgres + OpenSearch
- **Key Ops:** Search, Ingest, Articles  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-order-api

**Overview**  
Cart and checkout orchestration.

**Misc**  
- **Language/Stack:** Golang 
- **Port:** 7009  
- **Key Ops:** Cart, Checkout, Orders, Tax  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-order-discount-api

**Overview**  
Promotions and coupon rules.

**Misc**  
- **Language/Stack:** Golang
- **Port:** 7010  
- **Key Ops:** Evaluate, Manage rules  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-order-scheduling-api

**Overview**  
Delivery and pickup scheduling.

**Misc**  
- **Language/Stack:** Golang
- **Port:** 7011  
- **Key Ops:** Slots, Reserve, Cancel reservation  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-payments-api

**Overview**  
Payment methods, intents, and provider abstraction.

**Misc**  
- **Language/Stack:** Node + Fastify
- **Port:** 7012  
- **Key Ops:** Payment intents, Confirm, Webhooks  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-search-api

**Overview**  
Catalog search, autocomplete, ranking.

**Misc**  
- **Language/Stack:** Golang
- **Port:** 7013  
- **Data:** OpenSearch/Meilisearch  
- **Key Ops:** Search, Suggest, Reindex  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-servicing-chat-api

**Overview**  
Provides customer–agent chat integration for the e-commerce platform. Handles live support sessions, message delivery, and session lifecycle management.

**Misc**  
- **Language/Stack:** Node + Fastify
- **Port:** 7021  
- **Data:** OpenSearch/Meilisearch, Redis (session state)
- **Key Ops:** Message, Terminate 
- **Docs/Health:** `/docs`, `/health`

---

### komodo-ssr-engine

**Overview**  
Server-side rendering for storefront.

**Misc**  
- **Language/Stack:** Node + Next.js
- **Port:** 7014  
- **Key Ops:** Pages, Products, Categories  
- **Docs/Health:** `/health`

---

### komodo-user-api

**Overview**  
User profiles, addresses, preferences.

**Misc**  
- **Language/Stack:** Golang
- **Port:** 7016  
- **Key Ops:** User CRUD, Addresses  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-user-marketing-api

**Overview**  
Email/SMS campaigns and subscriptions.

**Misc**  
- **Language/Stack:** Node + Fastify
- **Port:** 7018  
- **Key Ops:** Subscribe, Campaigns, Consent  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-user-recommendations-api

**Overview**  
Personalized recommendations.

**Misc**  
- **Language/Stack:** Node + Fastify
- **Port:** 7019  
- **Data:** Feature store + ML models  
- **Key Ops:** Recommendations, Model reload  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-user-reviews-api

**Overview**  
Product reviews, ratings, moderation.

**Misc**  
- **Language/Stack:** Node + Fastify
- **Port:** 7020  
- **Key Ops:** Reviews, Ratings, Moderation  
- **Docs/Health:** `/docs`, `/health`

---

## Testing, Observability & Security

**Testing**  
- Unit, integration, and smoke tests per service.  

**Observability**  
- Local stack: Jaeger/Grafana/Loki.  
- Services expose `/metrics` (Prometheus).  

**Security**  
- JWT validation in all public endpoints.  
- Optional mTLS between internal services.  
- Secrets in `.env` locally, secret manager in prod.  
- Rate limiting at edge.

---

## CI/CD

- **CI:** Lint, tests, build, OpenAPI validation.  
- **CD:** Per-service deploys to Kubernetes/Helm.  
- **Versioning:** Independent or monorepo-wide semantic versions.

---

## Troubleshooting

- **Port conflict:** Change `PORT` in `.env`.  
- **DB migrations fail:** Reset and re-run migrations.  
- **Auth errors locally:** Enable `BYPASS_AUTH` for dev.

---

## License

MIT (or your preferred license). See `LICENSE`.
