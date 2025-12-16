# User API Quick Reference

## Profile Size Parameter Usage

### M2M (Service-to-Service) Examples

#### Get Basic Profile
```bash
curl -X POST http://localhost:7021/m2m/users/usr_123 \
  -H "Authorization: Bearer <service-token>" \
  -H "Content-Type: application/json" \
  -d '{"size": "basic"}'
```

**Response:**
```json
{
  "user_id": "usr_123",
  "first_name": "Sarah",
  "last_name": "Johnson",
  "avatar_url": "https://cdn.komodo.com/avatars/usr_123.jpg"
}
```

#### Get Minimal Profile
```bash
curl -X POST http://localhost:7021/m2m/users/usr_123 \
  -H "Authorization: Bearer <service-token>" \
  -H "Content-Type: application/json" \
  -d '{"size": "minimal"}'
```

**Response:**
```json
{
  "user_id": "usr_123",
  "email": "sarah.johnson@example.com",
  "phone": "+1-415-555-0142",
  "first_name": "Sarah",
  "last_name": "Johnson",
  "password_hash": "$2a$10$...",
  "avatar_url": "https://cdn.komodo.com/avatars/usr_123.jpg"
}
```

#### Get Full Profile
```bash
curl -X POST http://localhost:7021/m2m/users/usr_123 \
  -H "Authorization: Bearer <service-token>" \
  -H "Content-Type: application/json" \
  -d '{"size": "full"}'
```

**Response:**
```json
{
  "user_id": "usr_123",
  "username": "sarah.johnson",
  "email": "sarah.johnson@example.com",
  "phone": "+1-415-555-0142",
  "first_name": "Sarah",
  "middle_initial": "M",
  "last_name": "Johnson",
  "password_hash": "$2a$10$...",
  "address": {
    "address_id": "addr_456",
    "alias": "Home",
    "line1": "742 Evergreen Terrace",
    "line2": "Apt 3B",
    "line3": "",
    "city": "San Francisco",
    "state": "CA",
    "zip_code": "94102",
    "country": "USA"
  },
  "preferences": {
    "language": "en-US",
    "timezone": "America/Los_Angeles"
  },
  "metadata": {
    "created_at": "2023-08-15T14:32:18Z",
    "updated_at": "2025-12-04T18:45:22Z",
    "last_login": "2025-12-04T16:22:15Z",
    "email_verified": true,
    "mfa_enabled": true
  },
  "avatar_url": "https://cdn.komodo.com/avatars/usr_123.jpg"
}
```

---

### Client (Browser/Mobile) Examples

#### Get My Profile (Basic)
```bash
curl -X POST http://localhost:7021/users/me \
  -H "Cookie: session_id=<session-token>" \
  -H "X-CSRF-Token: <csrf-token>" \
  -H "Content-Type: application/json" \
  -d '{"size": "basic"}'
```

#### Get My Profile (Full) - For Account Settings
```bash
curl -X POST http://localhost:7021/users/me \
  -H "Cookie: session_id=<session-token>" \
  -H "X-CSRF-Token: <csrf-token>" \
  -H "Content-Type: application/json" \
  -d '{"size": "full"}'
```

#### Update My Profile
```bash
curl -X PATCH http://localhost:7021/users/me/profile \
  -H "Cookie: session_id=<session-token>" \
  -H "X-CSRF-Token: <csrf-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Sarah",
    "last_name": "Johnson-Smith",
    "phone": "+1-415-555-9999"
  }'
```

---

## When to Use Each Size

### `basic` - Lightweight Display
**Use for:**
- User avatars in comments/reviews
- Author names in blog posts
- Quick user mentions
- List views with many users

**Performance:** ~100 bytes per user

---

### `minimal` - Authentication & Basic Management
**Use for:**
- Login verification
- Password reset flows
- Email/phone validation
- Basic account info displays

**Performance:** ~200 bytes per user

---

### `full` - Complete Profile
**Use for:**
- Account settings page
- Checkout flow (needs address)
- Profile edit forms
- Admin user management
- Personalization features (needs preferences)

**Performance:** ~500-800 bytes per user

---

## Default Behavior

If no `size` parameter is provided, the API defaults to `"basic"`:

```bash
# These are equivalent:
POST /m2m/users/usr_123
POST /m2m/users/usr_123 -d '{"size": "basic"}'
```

---

## Error Responses

### Missing Authentication
```json
{
  "error": "missing authorization header",
  "error_code": "ERR_INVALID_TOKEN"
}
```

### Invalid Size Parameter
The API accepts any value and defaults to `"basic"` if invalid. Valid values:
- `"basic"`
- `"minimal"`
- `"full"`

---

## Performance Tips

1. **Always use the smallest size needed** - Reduces bandwidth and database load
2. **Cache basic profiles** - They change infrequently
3. **Batch requests when possible** - Future: `/m2m/users/batch` endpoint
4. **Use compression** - Enable gzip for full profiles

---

## Migration from Old Routes

### Old (GET with query param)
```bash
# ❌ Old way (no longer supported)
GET /m2m/users/usr_123?size=full
```

### New (POST with body)
```bash
# ✅ New way
POST /m2m/users/usr_123
Content-Type: application/json

{"size": "full"}
```

### Why the change?
- Request body is more flexible for future parameters
- Avoids caching issues with query params
- Consistent with modern API patterns
- Better for complex filtering in the future
