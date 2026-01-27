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
  - [komodo-auth-api](#komodo-auth-api)
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
cp .env
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

- **API Style:** JSON over HTTP; REST first; gRPC optional for internal API calls.  
- **Auth:** JWT access tokens via `komodo-auth-api`. Service-to-service via mTLS or signed service tokens.  
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

## Port Allocation Strategy

Services are allocated port ranges by business domain, with highest-priority/most-critical services closest to 7001.

> **Note:** Port 7000 is reserved due to conflicts with macOS system services (AFS/AirPlay/Control Center). All services start at 7001.

| Port Range | Domain | Priority | Description |
|------------|--------|----------|-------------|
| 7001–7019 | Infrastructure & Mesh | CRITICAL | Sidecars, Service Discovery, and Config management. |
| 7020–7049 | Identity & Security | CRITICAL | Auth, IAM, Secret management (Vault), and Token validation. |
| 7050–7149 | Core Business Logic | HIGH | Orders, Payments, and Inventory |
| 7150–7249 | User & Engagement | MEDIUM | Profiles, Social, Loyalty, and Communications. |
| 7250–7349 | Emerging Tech & AI | MEDIUM | GenAI, LLM Orchestration, and Vector DB interfaces. |
| 7350–7449 | Observability & Ops | LOW | Logging exporters, health checks, and logging agents. |

---

### 7020–7049: Identity & Security
- **7020** - `komodo-auth-api` - OAuth 2.0 flows
- **7021-7049** - Reserved for future auth services

---

### 7150–7249: User Services
- **7150** - `komodo-user-api` - Core user profiles, preferences
- **7151-7249** - Reserved for future user services
