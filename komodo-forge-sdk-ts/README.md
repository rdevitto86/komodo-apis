# @komodo/forge-apis-node

A TypeScript SDK providing shared utilities, types, and domain logic for Komodo applications. Organized with clear separation between backend (Node.js/server), frontend (browser/client), and shared (universal) code.

## Installation

For internal use, install via pnpm with a local path or private registry:

```bash
pnpm add @komodo/forge-apis-node
```

Or using a local path during development:

```bash
pnpm add ../komodo-forge-apis-node
```

## Architecture

The SDK is organized into three main categories:

- **`backend`** - Server-side only modules (Node.js, Express, AWS, etc.)
- **`frontend`** - Client-side only modules (browser APIs, UI utilities)
- **`shared`** - Universal modules (types, domain logic, utilities)

## Usage

### Backend Developers (Node.js/Server)

Import backend-specific modules for server-side functionality:

```typescript
// Import all backend modules
import { backend } from '@komodo/forge-apis-node';

backend.logging.runtime.logger.info('Server started');
backend.middleware.authenticate();

// Or import specific backend modules
import { aws, db, logging, middleware } from '@komodo/forge-apis-node/backend';

logging.runtime.logger.info('Hello from server');
await db.query('SELECT * FROM users');

// Import individual backend modules
import * as logging from '@komodo/forge-apis-node/backend/logging';
import * as aws from '@komodo/forge-apis-node/backend/aws';
```

### Frontend Developers (Browser/Client)

Import frontend-specific modules for client-side functionality:

```typescript
// Import all frontend modules
import { frontend } from '@komodo/forge-apis-node';

frontend.api.fetchUser();

// Or import specific frontend modules
import { api } from '@komodo/forge-apis-node/frontend';

api.fetchUser();
```

### Shared Code (Both Backend & Frontend)

Import shared modules that work in both environments:

```typescript
// Import all shared modules
import { shared } from '@komodo/forge-apis-node';

const user: shared.types.User = { id: '123', name: 'John' };
shared.domains.auth.validateToken(token);

// Or import specific shared modules
import { domains, crypto, utils } from '@komodo/forge-apis-node/shared';

domains.auth.validateToken(token);
crypto.hash(password);

// Import types directly
import type { Product, Service, MarketingContent } from '@komodo/forge-apis-node/shared/types';

const product: Product = {
  id: '1',
  slug: 'product-1',
  name: 'My Product',
  // ...
};
```

### Mixed Usage Example

```typescript
// Backend API handler
import { logging, db } from '@komodo/forge-apis-node/backend';
import { domains } from '@komodo/forge-apis-node/shared';
import type { Product } from '@komodo/forge-apis-node/shared/types';

export async function getProduct(id: string): Promise<Product> {
  logging.runtime.logger.info(`Fetching product ${id}`);
  
  const product = await db.query<Product>('SELECT * FROM products WHERE id = ?', [id]);
  
  return domains.products.transform(product);
}
```

```typescript
// Frontend component
import { api } from '@komodo/forge-apis-node/frontend';
import { domains } from '@komodo/forge-apis-node/shared';
import type { Product } from '@komodo/forge-apis-node/shared/types';

export function ProductCard({ productId }: { productId: string }) {
  const [product, setProduct] = useState<Product | null>(null);
  
  useEffect(() => {
    api.fetchProduct(productId).then(setProduct);
  }, [productId]);
  
  return <div>{product?.name}</div>;
}
```

## Available Modules

### Backend Modules (Server-side only)

- **`backend/aws`** - AWS SDK utilities and helpers
- **`backend/config`** - Server configuration management
- **`backend/db`** - Database utilities and helpers
- **`backend/logging`** - Server-side logging
  - `runtime` - Runtime logging
  - `security` - Security logging
  - `telemetry` - Telemetry logging
- **`backend/middleware`** - Express/API middleware
- **`backend/observability`** - Server observability and monitoring

### Frontend Modules (Client-side only)

- **`frontend/api`** - API client wrappers and utilities

### Shared Modules (Universal)

- **`shared/crypto`** - Cryptography utilities
- **`shared/domains`** - Domain-specific business logic
  - `auth` - Authentication domain
  - `payments` - Payments domain
  - `user` - User domain
- **`shared/entitlements`** - Entitlement management
- **`shared/feature-flags`** - Feature flag utilities
- **`shared/security`** - Security utilities
- **`shared/utils`** - General utilities
- **`shared/types`** - TypeScript type definitions
  - `Product` - Product types
  - `Service` - Service types
  - `MarketingContent` - Marketing content types
  - `Campaign` - Campaign types
  - And more...

## Development

### Building

```bash
pnpm run build
```

### Cleaning build artifacts

```bash
pnpm run clean
```

### Rebuilding from scratch

```bash
pnpm run rebuild
```

## TypeScript Support

This library is written in TypeScript and includes full type definitions. TypeScript projects will automatically get IntelliSense and type checking when using this library.

## Project Structure

```
src/
├── backend/              # Server-side only
│   ├── aws/              # AWS utilities
│   ├── config/           # Configuration
│   ├── db/               # Database utilities
│   ├── logging/          # Logging
│   │   ├── runtime/
│   │   ├── security/
│   │   └── telemetry/
│   ├── middleware/       # Middleware
│   ├── observability/    # Observability
│   └── index.ts
├── frontend/             # Client-side only
│   ├── api/              # API clients
│   └── index.ts
├── shared/               # Universal (both backend & frontend)
│   ├── crypto/           # Cryptography
│   ├── domains/          # Domain logic
│   │   ├── auth/
│   │   ├── payments/
│   │   └── user/
│   ├── entitlements/     # Entitlements
│   ├── feature-flags/    # Feature flags
│   ├── security/         # Security
│   ├── types/            # Type definitions
│   │   ├── marketing.d.ts
│   │   ├── orders.d.ts
│   │   ├── products.d.ts
│   │   └── services.d.ts
│   ├── utils/            # Utilities
│   └── index.ts
└── index.ts              # Main entry point
```

## Publishing to Private NPM Registry

When ready to publish to your private npm registry:

1. Update version in `package.json`
2. Build the package: `pnpm run build`
3. Publish: `pnpm publish --registry=<your-private-registry-url>`

## License

ISC - Internal use only
