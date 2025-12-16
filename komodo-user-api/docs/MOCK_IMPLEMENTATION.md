# Mock Implementation Summary

## Overview

All handlers now load mock data from JSON files during development, with commented-out production DynamoDB implementations ready to be enabled.

---

## Mock Files

Located in `/mocks/`:

### `user_basic.json`
```json
{
  "user_id": "usr_2k8f9j3h4g5d6s7a",
  "first_name": "Sarah",
  "last_name": "Johnson",
  "avatar_url": "https://cdn.komodo.com/avatars/usr_2k8f9j3h4g5d6s7a.jpg"
}
```

### `user_minimal.json`
```json
{
  "user_id": "usr_7m2n9p4q8r5t3v1w",
  "email": "michael.chen@example.com",
  "phone": "+1-415-555-0142",
  "first_name": "Michael",
  "last_name": "Chen",
  "password_hash": "$2a$10$...",
  "avatar_url": "https://cdn.komodo.com/avatars/usr_7m2n9p4q8r5t3v1w.jpg"
}
```

### `user_full.json`
```json
{
  "user_id": "usr_5x8c2v9b4n7m3k1j",
  "username": "alexandra.rodriguez",
  "email": "alexandra.rodriguez@example.com",
  "phone": "+1-650-555-0198",
  "first_name": "Alexandra",
  "middle_initial": "M",
  "last_name": "Rodriguez",
  "password_hash": "$2a$10$...",
  "address": { ... },
  "preferences": { ... },
  "metadata": { ... },
  "avatar_url": "https://cdn.komodo.com/avatars/usr_5x8c2v9b4n7m3k1j.jpg"
}
```

---

## Handler Implementations

### M2M Handlers (`internal/handlers/m2m_handlers.go`)

#### `GetUserByID`
- **Mock**: Loads JSON file based on `size` parameter (basic/minimal/full)
- **Production (commented)**: DynamoDB `GetItem` with projection expression
- **Features**:
  - Validates size parameter
  - Overrides `user_id` with request parameter
  - Returns appropriate fields based on size

#### `CreateUser`
- **Mock**: Returns success with submitted data
- **Production (commented)**: DynamoDB `PutItem` with validation

#### `UpdateUserByID`
- **Mock**: Returns success with update data
- **Production (commented)**: DynamoDB `UpdateItem` with expression builder

#### `DeleteUserByID`
- **Mock**: Returns success message
- **Production (commented)**: DynamoDB soft delete (sets `account_status` and `deleted_at`)

---

### User Profile Handlers (`internal/handlers/user_profile_handlers.go`)

#### `GetMyProfile`
- **Mock**: Loads JSON file based on `size` parameter
- **Production (commented)**: Same as `GetUserByID` but uses session `user_id`
- **Features**:
  - Validates size parameter
  - Overrides `user_id` with session value
  - Returns appropriate fields based on size

#### `UpdateMyProfile`
- **Mock**: Returns success with update data
- **Production (commented)**: DynamoDB `UpdateItem` with dynamic expression builder
  - Builds update expression from request body
  - Automatically adds `updated_at` timestamp
  - Uses expression attribute names/values for safety

#### `DeleteMyAccount`
- **Mock**: Returns success message
- **Production (commented)**: DynamoDB soft delete
  - Sets `account_status = "deleted"`
  - Sets `deleted_at` timestamp

---

### Address Handlers (`internal/handlers/user_addresses_handlers.go`)

#### `GetMyAddresses`
- **Mock**: Returns array with one sample address
- **Production (commented)**: DynamoDB projection of `addresses` nested attribute
- **Data Model**: Addresses stored as nested array in user record

#### `AddMyAddress`
- **Mock**: Returns success with submitted address
- **Production (commented)**: DynamoDB `UpdateItem` to append to `addresses` list

#### `UpdateMyAddress`
- **Mock**: Returns success with update data
- **Production (commented)**: DynamoDB `UpdateItem` to modify specific address in list

#### `DeleteMyAddress`
- **Mock**: Returns success message
- **Production (commented)**: DynamoDB `UpdateItem` to remove address from list

---

### Preference Handlers (`internal/handlers/user_preferences_handlers.go`)

#### `GetMyPreferences`
- **Mock**: Returns `language` and `timezone`
- **Production (commented)**: DynamoDB projection of `preferences` nested attribute
- **Data Model**: Preferences stored as nested map in user record

#### `UpdateMyPreferences`
- **Mock**: Returns success with submitted preferences
- **Production (commented)**: DynamoDB `UpdateItem` to modify `preferences` map

---

## DynamoDB Data Model

### User Table Schema

**Table Name**: `komodo-users`

**Primary Key**: `user_id` (String)

**Attributes**:
```
user_id (String, PK)
username (String)
email (String)
phone (String)
first_name (String)
middle_initial (String)
last_name (String)
password_hash (String)
avatar_url (String)

address (Map)
  ├─ address_id (String)
  ├─ alias (String)
  ├─ line1 (String)
  ├─ line2 (String)
  ├─ line3 (String)
  ├─ city (String)
  ├─ state (String)
  ├─ zip_code (String)
  └─ country (String)

addresses (List of Maps) - For multiple addresses
  └─ [Same structure as address]

preferences (Map)
  ├─ language (String)
  └─ timezone (String)

metadata (Map)
  ├─ created_at (String, ISO 8601)
  ├─ updated_at (String, ISO 8601)
  ├─ last_login (String, ISO 8601)
  ├─ email_verified (Boolean)
  └─ mfa_enabled (Boolean)

account_status (String) - "active" | "deleted" | "suspended"
deleted_at (String, ISO 8601) - Set on soft delete
```

---

## Switching to Production

To enable DynamoDB:

1. **Uncomment the DynamoDB code** in each handler
2. **Add DynamoDB client initialization** in `main.go`:
   ```go
   import (
       "github.com/aws/aws-sdk-go-v2/config"
       "github.com/aws/aws-sdk-go-v2/service/dynamodb"
   )
   
   cfg, _ := config.LoadDefaultConfig(context.Background())
   dynamoClient := dynamodb.NewFromConfig(cfg)
   ```
3. **Set environment variable**:
   ```bash
   export DYNAMODB_USERS_TABLE=komodo-users
   ```
4. **Create DynamoDB table** with the schema above
5. **Remove or conditionally disable** mock file loading

---

## Benefits of This Approach

✅ **Development Ready**: Works immediately with mock data  
✅ **Production Ready**: DynamoDB code is complete, just commented out  
✅ **Easy Testing**: Mock files can be modified for different test scenarios  
✅ **Clear Migration Path**: Uncomment blocks to switch to production  
✅ **Performance**: Size parameter reduces data transfer  
✅ **Flexibility**: Nested attributes keep related data together  

---

## Testing the Mocks

### Start the server:
```bash
cd /Users/rad/komodo-apis/komodo-user-api
./build/komodo-user-api
```

### Test M2M endpoint:
```bash
curl -X POST http://localhost:7021/m2m/users/test123 \
  -H "Authorization: Bearer mock-token" \
  -H "Content-Type: application/json" \
  -d '{"size": "full"}'
```

### Expected Response:
Full user profile from `user_full.json` with `user_id` overridden to `test123`

---

## Next Steps

1. **Create DynamoDB table** in AWS
2. **Implement DynamoDB client singleton** in shared library
3. **Add data validation** before database writes
4. **Implement proper error handling** for database operations
5. **Add logging** for database queries
6. **Create database migration scripts** for schema changes
7. **Add indexes** for common query patterns (email, username lookups)
