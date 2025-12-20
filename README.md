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
| 7001–7020 | Authentication & Authorization | CRITICAL | Core auth services |
| 7021–7040 | User Services | HIGH | User profiles & preferences |
| 7041–7080 | Order & Payments | HIGH | Transaction services |
| 7081–7120 | Catalog & Inventory | HIGH | Product discovery |
| 7121–7160 | Core Platform | MEDIUM | Platform infrastructure |
| 7161–7200 | Customer Engagement | MEDIUM | Support & communications |
| 7201–7220 | Reviews & Loyalty | MEDIUM | Social & retention |
| 7221–7250 | GenAI & Automation | LOW | AI features |
| 7251–7290 | Analytics & Observability | LOW | Monitoring & metrics |
| 7291–7300 | SSR & Rendering | LOW | Server-side rendering |

---

### 7001–7020: Authentication & Authorization
**Priority: CRITICAL** - Core authentication/authorization services that all other services depend on
- **7001** - `komodo-auth-api` - OAuth 2.0 flows
- **7002-7020** - Reserved for future auth services

---

### 7021–7040: User Services
**Priority: HIGH** - User profile and account management

- **7021** - `komodo-user-api` - Core user profiles, preferences
- **7022** - `komodo-user-personalize-api` - Personalization, recommendations
- **7023-7040** - Reserved for future user services

---

### 7041–7080: Order & Payments
**Priority: HIGH** - Revenue-generating transaction services

#### Order Services (7041–7060)
- **7041** - `komodo-order-manage-api` - Create, update orders
- **7042** - `komodo-order-lookup-api` - Order history, tracking
- **7043** - `komodo-order-discounts-api` - Promo codes, discounts
- **7044** - `komodo-order-returns-api` - Returns, refunds
- **7045** - `komodo-order-shipping-api` - Shipping calculations, tracking
- **7046-7060** - Reserved for future order services

#### Payment Services (7061–7080)
- **7061** - `komodo-payments-intents-api` - Payment intents, checkout
- **7062** - `komodo-payments-transactions-api` - Transaction processing
- **7063** - `komodo-payments-methods-api` - Stored payment methods
- **7064** - `komodo-payments-ledger-api` - Financial ledger, reconciliation
- **7065-7080** - Reserved for future payment services

---

### 7081–7120: Catalog & Inventory
**Priority: HIGH** - Product discovery and availability

#### Catalog Services (7081–7100)
- **7081** - `komodo-catalog-items-read-api` - Product listings, details
- **7082** - `komodo-catalog-items-manage-api` - Product CRUD
- **7083** - `komodo-catalog-items-suggestion-api` - AI-powered suggestions
- **7084** - `komodo-catalog-inventory-read-api` - Stock levels
- **7085** - `komodo-catalog-inventory-manage-api` - Inventory updates
- **7086-7100** - Reserved for future catalog services

#### Search Services (7101–7120)
- **7101** - `komodo-search-api` - Product search, filtering
- **7102** - `komodo-search-suggestions-api` - Autocomplete, typeahead
- **7103-7120** - Reserved for future search services

---

### 7121–7160: Core Platform Services
**Priority: MEDIUM** - Essential platform capabilities

#### Entitlements & Features (7121–7130)
- **7121** - `komodo-core-entitlements-api` - User permissions, access control
- **7122** - `komodo-core-features-api` - Feature flags, A/B testing
- **7123** - `komodo-core-events-api` - Event bus, pub/sub
- **7124-7130** - Reserved for future core services

#### Content & Media (7131–7140)
- **7131** - `komodo-content-cdn-api` - CDN management, media delivery
- **7132-7140** - Reserved for future content services

#### Address & Location (7141–7150)
- **7141** - `komodo-address-api` - Address validation, geocoding
- **7142-7150** - Reserved for future location services

---

### 7161–7200: Customer Engagement
**Priority: MEDIUM** - Customer support and communication

#### Support Services (7161–7180)
- **7161** - `komodo-support-cases-api` - Support tickets, cases
- **7162** - `komodo-support-agent-chat-api` - Live agent chat
- **7163** - `komodo-support-knowledge-api` - KB articles, FAQs
- **7164** - `komodo-support-feedback-api` - Customer feedback, surveys
- **7165-7180** - Reserved for future support services

#### Communications (7181–7200)
- **7181** - `komodo-comms-messaging-api` - Email, SMS, push notifications
- **7182** - `komodo-comms-preferences-api` - Notification preferences
- **7183-7200** - Reserved for future comms services

---

### 7201–7220: Reviews & Loyalty
**Priority: MEDIUM** - Social proof and retention

#### Reviews (7201–7210)
- **7201** - `komodo-reviews-read-api` - Read reviews, ratings
- **7202** - `komodo-reviews-submit-api` - Submit reviews
- **7203** - `komodo-reviews-moderation-api` - Moderation, flagging
- **7204-7210** - Reserved for future review services

#### Loyalty (7211–7220)
- **7211** - `komodo-loyalty-rewards-api` - Points, rewards, tiers
- **7212-7220** - Reserved for future loyalty services

---

### 7221–7250: GenAI & Automation
**Priority: LOW** - AI-powered features and automation

- **7221** - `komodo-genai-agent-chat-api` - AI chatbot, customer service
- **7222** - `komodo-genai-summary-api` - AI summaries, descriptions
- **7223-7250** - Reserved for future AI services

---

### 7251–7290: Analytics & Observability
**Priority: LOW** - Monitoring, metrics, and analytics

- **7251** - `komodo-analytics-interaction-api` - User interactions, clickstream
- **7252** - `komodo-analytics-logs-api` - Centralized logging
- **7253** - `komodo-analytics-telemetry-api` - Metrics, traces
- **7254-7290** - Reserved for future analytics services

---

### 7291–7300: SSR & Rendering
**Priority: LOW** - Server-side rendering engines

- **7291** - `komodo-ssr-engine` - React/Next.js SSR
- **7292-7300** - Reserved for future rendering services
