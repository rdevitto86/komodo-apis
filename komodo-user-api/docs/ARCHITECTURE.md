# User API Architecture

## Route Structure

The User API maintains **two separate route groups** for different client types:

### 1. M2M Routes (Service-to-Service)
**Base Path:** `/m2m`

**Authentication:** JWT or OAuth service tokens (Bearer tokens only)

**Middleware Stack:**
- `ServiceAuthMiddleware()` - Validates service tokens

**Security Characteristics:**
- No CSRF protection needed (not browser-based)
- Higher rate limits (trusted internal services)
- No session management
- No idempotency checks

**Endpoints:**
```
POST   /m2m/users          - Create new user
POST   /m2m/users/:id      - Get user by ID (with size param)
PATCH  /m2m/users/:id      - Update user by ID
DELETE /m2m/users/:id      - Delete user by ID
```

---

### 2. Client Routes (Browser/Mobile Apps)
**Base Path:** `/users/me`

**Authentication:** Session-based (cookies) + CSRF tokens

**Middleware Stack (when enabled):**
- `Session()` - Validates session cookies
- `CSRF()` - Protects against cross-site request forgery
- `Idempotency()` - Prevents duplicate mutations
- `RateLimiter()` - Browser-specific rate limits

**Security Characteristics:**
- Full browser security protections
- CSRF validation required for mutations
- Session-based authentication
- Idempotency for safe retries

**Endpoints:**
```
POST   /users/me                    - Get my profile (with size param)
PATCH  /users/me/profile            - Update my profile
DELETE /users/me/account            - Delete my account

GET    /users/me/addresses          - List my addresses
POST   /users/me/addresses          - Add new address
PATCH  /users/me/addresses/:addr_id - Update address
DELETE /users/me/addresses/:addr_id - Delete address

GET    /users/me/preferences        - Get my preferences
PATCH  /users/me/preferences        - Update my preferences
```

---

## Why Separate Routes?

### ❌ Why NOT Auto-Detect Client Type?

We considered using a unified middleware that auto-detects Bearer vs Session auth, but decided against it for these reasons:

1. **Security Boundaries**
   - M2M and Client routes have fundamentally different security requirements
   - CSRF protection is critical for browsers but unnecessary for services
   - Mixing auth types creates potential bypass vulnerabilities

2. **Different Middleware Stacks**
   - Services need minimal overhead (just token validation)
   - Browsers need full protection (CSRF, sessions, idempotency, rate limiting)
   - Conditional middleware based on auth type adds complexity

3. **Clear Separation of Concerns**
   - Internal service traffic vs external user traffic
   - Different SLAs and monitoring requirements
   - Easier to audit and debug

4. **Prevents Security Bypasses**
   - Browser clients can't bypass CSRF by sending Bearer tokens
   - Services can't accidentally trigger session logic

### ✅ Advantages of Separate Routes

1. **Explicit Security Policies** - Each route group has clear, non-negotiable security requirements
2. **Better Observability** - Logs and metrics clearly distinguish internal vs external traffic
3. **Flexible Rate Limiting** - Different limits for trusted services vs public clients
4. **Simpler Middleware** - No conditional logic based on auth type
5. **Future-Proof** - Easy to add route-group-specific features (e.g., service-specific quotas)

---

## Profile Size Parameter

Both M2M and Client routes support a `size` parameter to control how much user data is returned:

### Size Options

| Size | Fields Included | Use Case |
|------|----------------|----------|
| `basic` | `user_id`, `first_name`, `last_name`, `avatar_url` | Lightweight display (comments, reviews) |
| `minimal` | Basic + `email`, `phone`, `password_hash` | Authentication, basic account management |
| `full` | Minimal + `username`, `middle_initial`, `address`, `preferences`, `metadata` | Complete profile, checkout, account settings |

### Request Format

**M2M Route:**
```bash
POST /m2m/users/usr_123
Content-Type: application/json
Authorization: Bearer <service-token>

{
  "size": "full"
}
```

**Client Route:**
```bash
POST /users/me
Content-Type: application/json
Cookie: session_id=<session-token>
X-CSRF-Token: <csrf-token>

{
  "size": "minimal"
}
```

### Why POST Instead of GET?

We use `POST` for profile retrieval (instead of `GET`) to:
1. **Accept request body** - Size parameter in body is cleaner than query params
2. **Avoid caching issues** - Different sizes shouldn't be cached as the same resource
3. **Future extensibility** - Easy to add filters, projections, or other parameters
4. **Consistent with GraphQL patterns** - Similar to GraphQL queries with variable depth

---

## Middleware Implementation Status

| Middleware | Status | Notes |
|------------|--------|-------|
| `ServiceAuthMiddleware` | ✅ Implemented | JWT + OAuth validation |
| `Session` | 🚧 Partial | Chi version exists, Gin version needed |
| `CSRF` | 🚧 Planned | Requires session middleware first |
| `Idempotency` | 🚧 Planned | For mutation safety |
| `RateLimiter` | 🚧 Planned | Different limits per route group |

---

## Data Structures

See `komodo-forge-sdk-go/domains/user/user_types.go` for complete type definitions:

- `UserProfileGetResponseBasic` - Basic profile fields
- `UserProfileGetResponseMinimal` - Minimal profile fields
- `UserProfileGetResponseFull` - Full profile with sub-objects:
  - `UserAddress` - Address information
  - `UserPreferences` - Language, timezone
  - `UserMetadata` - Timestamps, verification status

---

## Future Considerations

1. **GraphQL Alternative** - Consider GraphQL for flexible field selection
2. **Field Projection** - Allow clients to specify exact fields needed
3. **Caching Strategy** - Different cache TTLs per size level
4. **Compression** - Compress full profiles for bandwidth savings
