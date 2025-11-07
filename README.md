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
  - [komodo-auth-internal-api](#komodo-auth-internal-api)
  - [komodo-address-api](#komodo-address-api)
  - [komodo-ai-chatbot-api](#komodo-ai-chatbot-api)
  - [komodo-ai-summary-api](#komodo-ai-summary-api)
  - [komodo-analytics-interaction-api](#komodo-analytics-interaction-api)
  - [komodo-analytics-logs-api](#komodo-analytics-logs-api)
  - [komodo-analytics-telemetry-api](#komodo-analytics-telemetry-api)
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
- **Auth:** JWT access tokens via `komodo-auth-internal-api`. Service-to-service via mTLS or signed service tokens.  
- **Tracing:** OpenTelemetry (OTLP) → collector (Jaeger/Tempo).  
- **Logging:** Structured JSON, request IDs, 12-factor log to stdout.  
- **Errors:** RFC 7807 (Problem+JSON) recommended.  
- **Schema:** OpenAPI per service at `/docs` and `/openapi.json`.

---

## Services

Each service uses the template:

- **Overview**  
- **Build**  
- **Run**  
- **Misc** (Language/Stack, Port, Data Stores, Key Operations, Docs, Health)  

---

### Core Services (7000–7099)

---

### komodo-auth-internal-api

**Overview**  
Authentication, authorization, MFA, etc.

**Misc**  
- **Language/Stack:** Golang
- **Port:** DEV: 7001 | PROD: 8080
- **Key Ops:** Signup, Login, Token refresh, Password reset  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-entitlements-api

**Overview**  
Manages user entitlements and access control.

**Misc**  
- **Language/Stack:** Golang
- **Port:** DEV: 7002 | PROD: 8080
- **Key Ops:** Entitlement checks, Role management  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-analytics-interaction-api

**Overview**  
Captures user interactions for real-time personalization.

**Misc**  
- **Language/Stack:** Golang
- **Port:** DEV: 7003 | PROD: 8080
- **Data:** Kafka/Redpanda  
- **Key Ops:** Events ingestion, Stats query  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-analytics-logs-api

**Overview**  
Ingests application logs and exposes query endpoints.

**Misc**  
- **Language/Stack:** Golang
- **Port:** DEV: 7004 | PROD: 8080
- **Data:** Loki/Elastic  
- **Key Ops:** Ingest, Query logs  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-analytics-telemetry-api

**Overview**  
Metrics/trace ingestion + query proxy.

**Misc**  
- **Language/Stack:** Golang
- **Port:** DEV: 7005 | PROD: 8080
- **Data:** Prometheus, Jaeger/Tempo  
- **Key Ops:** Metrics, Traces, Dashboards  
- **Docs/Health:** `/docs`, `/health`

---

### Address APIs (7100–7199)

---

### komodo-address-api

**Overview**  
Validates, normalizes, and geocodes shipping/billing addresses.

**Misc**  
- **Language/Stack:** Golang
- **Port:** DEV: 7100 | PROD: 8080
- **Data Stores:** Postgres, Redis
- **Key Ops:** Validate, Normalize, Geocode  
- **Docs/Health:** `/docs`, `/health`

---

### AI APIs (7200–7299)

---

### komodo-ai-chatbot-api

**Overview**  
Provides LLM-generated servicing via self-service chat.

**Misc**  
- **Language/Stack:** Python + FastAPI
- **Port:** DEV: 7200 | PROD: 8080
- **Upstreams:** LLM provider(s)
- **Key Ops:** Message  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-ai-summary-api

**Overview**  
Generates product/order/user summaries (LLMs).

**Misc**  
- **Language/Stack:** Python + FastAPI
- **Port:** DEV: 7201 | PROD: 8080
- **Upstreams:** LLM provider(s)  
- **Key Ops:** Summaries, Models listing  
- **Docs/Health:** `/docs`, `/health`

---

### Catalog APIs (7300–7399)

---

### komodo-catalog-api

**Overview**  
Products, variants, categories, pricing, inventory.

**Misc**  
- **Language/Stack:** Golang
- **Port:** DEV: 7300 | PROD: 8080
- **Data:** Postgres, Redis, S3  
- **Key Ops:** Products CRUD, Inventory, Categories  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-catalog-engagement-api

**Overview**  
Tracks and analyzes catalog engagement metrics.

**Misc**  
- **Language/Stack:** Golang
- **Port:** DEV: 7301 | PROD: 8080
- **Key Ops:** Engagement tracking, Metrics  
- **Docs/Health:** `/docs`, `/health`

---

### Communication APIs (7400–7499)

---

### komodo-comms-api

**Overview**  
Handles email, SMS, and push notifications.

**Misc**  
- **Language/Stack:** Node.js
- **Port:** DEV: 7400 | PROD: 8080
- **Key Ops:** Send notifications, Manage templates  
- **Docs/Health:** `/docs`, `/health`

---

### Order APIs (7500–7599)

---

### komodo-order-api

**Overview**  
Cart and checkout orchestration.

**Misc**  
- **Language/Stack:** Golang
- **Port:** DEV: 7500 | PROD: 8080
- **Key Ops:** Cart, Checkout, Orders, Tax  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-order-discount-api

**Overview**  
Promotions and coupon rules.

**Misc**  
- **Language/Stack:** Golang
- **Port:** DEV: 7501 | PROD: 8080
- **Key Ops:** Evaluate, Manage rules  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-order-returns-api

**Overview**  
Handles order returns and refunds.

**Misc**  
- **Language/Stack:** Golang
- **Port:** DEV: 7502 | PROD: 8080
- **Key Ops:** Returns, Refunds  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-order-scheduling-api

**Overview**  
Delivery and pickup scheduling.

**Misc**  
- **Language/Stack:** Golang
- **Port:** DEV: 7503 | PROD: 8080
- **Key Ops:** Slots, Reserve, Cancel reservation  
- **Docs/Health:** `/docs`, `/health`

---

### Payments APIs (7600–7699)

---

### komodo-payments-api

**Overview**  
Payment methods, intents, and provider abstraction.

**Misc**  
- **Language/Stack:** Node.js
- **Port:** DEV: 7600 | PROD: 8080
- **Key Ops:** Payment intents, Confirm, Webhooks  
- **Docs/Health:** `/docs`, `/health`

---

### Search APIs (7700–7799)

---

### komodo-search-api

**Overview**  
Catalog search, autocomplete, ranking.

**Misc**  
- **Language/Stack:** Golang
- **Port:** DEV: 7700 | PROD: 8080
- **Data:** OpenSearch/Meilisearch  
- **Key Ops:** Search, Suggest, Reindex  
- **Docs/Health:** `/docs`, `/health`

---

### User APIs (7800–7899)

---

### komodo-user-api

**Overview**  
User profiles, addresses, preferences.

**Misc**  
- **Language/Stack:** Golang
- **Port:** DEV: 7800 | PROD: 8080
- **Key Ops:** User CRUD, Addresses  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-user-marketing-api

**Overview**  
Email/SMS campaigns and subscriptions.

**Misc**  
- **Language/Stack:** Node.js
- **Port:** DEV: 7801 | PROD: 8080
- **Key Ops:** Subscribe, Campaigns, Consent  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-user-recommendations-api

**Overview**  
Personalized recommendations.

**Misc**  
- **Language/Stack:** Node.js
- **Port:** DEV: 7802 | PROD: 8080
- **Data:** Feature store + ML models  
- **Key Ops:** Recommendations, Model reload  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-user-reviews-api

**Overview**  
Product reviews, ratings, moderation.

**Misc**  
- **Language/Stack:** Node.js
- **Port:** DEV: 7803 | PROD: 8080
- **Key Ops:** Reviews, Ratings, Moderation  
- **Docs/Health:** `/docs`, `/health`

---

### komodo-user-rewards-api

**Overview**  
Manages user rewards and loyalty programs.

**Misc**  
- **Language/Stack:** Node.js
- **Port:** DEV: 7804 | PROD: 8080
- **Key Ops:** Rewards, Loyalty programs  
- **Docs/Health:** `/docs`, `/health`
