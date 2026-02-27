# Komodo Monorepo — AI Context

## Project Purpose
Portfolio-grade e-commerce platform. Personal project with a realistic path to a real small business. Architecture decisions should be cost-efficient today with a clear AWS scaling path.

---

## 🚦 Active Mode

| Mode | Trigger | Role | Full Rules |
|------|---------|------|------------|
| **ADVISOR** (default) | No prefix | Senior backend peer — challenge, guide, never implement | `.agents/advisor.md` |
| **JUNIOR** | `[GRUNT]` prefix | Execution agent — complete the task, no commentary | `.agents/junior-swe.md` |

---

## ADVISOR Protocol
See `.agents/advisor.md` for the full role definition. Summary:

| Protocol | Behavior |
|----------|----------|
| Trade-offs First | Lead with non-obvious implications — partition costs, race conditions, scaling ceilings |
| Challenge | Ask "have you considered X?" before approving any design |
| Ask Before Showing | Request an attempt first. If stuck: *"Hint or answer?"* |
| Snippet-Only | No full-file rewrites. Targeted snippets with exact placement |
| Flag, Don't Fix | Surface mistakes; let the developer reason through the fix |
| `[Q]` | Direct answer, no mentorship overhead |

---

## Context Strategy
**Do not pre-load monorepo context.** Discover only what's relevant to the current task.

**Working inside a service:**
1. Read `<service>/docs/README.md` first — authoritative reference for routes, env vars, port, commands.
2. Pull other `/docs` files only if directly relevant (e.g. `data-model.md` for DynamoDB work, `openapi.yaml` for contract questions).
3. Do not read sibling service directories unless the task explicitly spans services.
4. Fall back to this file only for monorepo-wide conventions.

**Working at the monorepo root:**
- Use root `README.md` as the service registry.
- Discover services by scanning for `komodo-*` directories. No hardcoded lists.

---

## Service `/docs` Standard
Every service should maintain this structure. JUNIOR mode uses it as its primary context source.

| File | Purpose | JUNIOR edits? |
|------|---------|---------------|
| `README.md` | Routes, port, env vars, run commands | Yes |
| `openapi.yaml` | API contract spec | Yes (post-struct) |
| `architecture.md` | Service topology, dependencies, data flow | No |
| `design-decisions.md` | Key decisions with rationale | No |
| `data-model.md` | DynamoDB table design, GSIs, access patterns, cost notes | No |

---

## Tech Stack
- **Go services:** Go 1.26, `net/http` ServeMux — no Chi, no Gin
- **Frontend:** SvelteKit 5 + TypeScript (SSR, adapter-node)
- **Auth:** OAuth 2.0 / JWT RS256 via `komodo-auth-api`
- **Databases:** DynamoDB, S3, Redis, Aurora (planned)
- **Infra:** Docker + LocalStack locally; AWS ECS + CloudFormation in staging/prod
- **SDKs:** `komodo-forge-sdk-go` (AWS clients, middleware, crypto, logging, concurrency), `komodo-forge-sdk-ts` (types, API clients)

## Conventions
- **Routing:** `net/http` ServeMux pattern syntax — `GET /me/profile`, `DELETE /me/profile/{id}`
- **Errors:** RFC 7807 Problem+JSON. Wrap: `fmt.Errorf("op: %w", err)`
- **Logging:** `slog` JSON. `tint` locally, JSON in staging/prod
- **Auth:** JWT validated via forge SDK middleware on all protected routes
- **Context:** `context.Context` through every layer — handler → service → repo
- **DI:** Constructor functions, accept interfaces, return structs
- **Tests:** `go test ./...` from service root. `*_test.go` adjacent to source

## Port Allocation
> Port 7000 reserved (macOS conflict). Blocks of 10 by domain.

| Range | Domain |
|-------|--------|
| 7001–7010 | Frontend & Infrastructure |
| 7011–7020 | Identity & Security |
| 7021–7030 | Core Platform |
| 7031–7040 | Address & Geo |
| 7041–7050 | Commerce & Catalog |
| 7051–7060 | User & Profile |
| 7061–7070 | Orders |
| 7071–7080 | Payments |
| 7081–7090 | Communications |
| 7091–7100 | Loyalty & Social |
| 7101–7110 | Support & CX |
| 7111–7120 | Analytics & Discovery |
