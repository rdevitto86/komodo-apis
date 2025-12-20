# @komodo/forge-apis-node

An internal TypeScript library providing shared utilities, middleware, and domain logic for Komodo Node.js APIs.

## Installation

For internal use, install via npm with a local path or private registry:

```bash
npm install @komodo/forge-apis-node
```

Or using a local path during development:

```bash
npm install ../komodo-forge-apis-node
```

## Usage

### Importing the entire library

```typescript
import * as forge from '@komodo/forge-apis-node';

// Use namespaced modules
forge.logging.runtime.logger.info('Hello');
forge.utils.someUtility();
```

### Importing specific modules

```typescript
import { logging, utils, middleware } from '@komodo/forge-apis-node';

// Use directly
logging.runtime.logger.info('Hello');
utils.someUtility();
```

### Importing from subpaths

```typescript
import * as logging from '@komodo/forge-apis-node/logging';
import * as utils from '@komodo/forge-apis-node/utils';
import * as middleware from '@komodo/forge-apis-node/middleware';

// Use the imported modules
logging.runtime.logger.info('Hello');
```

## Available Modules

- **aws** - AWS SDK utilities and helpers
- **config** - Configuration management
- **crypto** - Cryptography utilities
- **db** - Database utilities and helpers
- **domains** - Domain-specific business logic
  - `auth` - Authentication domain
  - `payments` - Payments domain
  - `user` - User domain
- **entitlements** - Entitlement management
- **feature-flags** - Feature flag utilities
- **logging** - Logging utilities
  - `runtime` - Runtime logging
  - `security` - Security logging
  - `telemetry` - Telemetry logging
- **middleware** - Express/API middleware
- **observability** - Observability and monitoring
- **security** - Security utilities
- **utils** - General utilities

## Development

### Building

```bash
npm run build
```

### Cleaning build artifacts

```bash
npm run clean
```

### Rebuilding from scratch

```bash
npm run rebuild
```

## TypeScript Support

This library is written in TypeScript and includes full type definitions. TypeScript projects will automatically get IntelliSense and type checking when using this library.

## Project Structure

```
src/
├── aws/              # AWS utilities
├── config/           # Configuration
├── crypto/           # Cryptography
├── db/               # Database utilities
├── domains/          # Domain logic
│   ├── auth/
│   ├── payments/
│   └── user/
├── entitlements/     # Entitlements
├── feature-flags/    # Feature flags
├── logging/          # Logging
│   ├── runtime/
│   ├── security/
│   └── telemetry/
├── middleware/       # Middleware
├── observability/    # Observability
├── security/         # Security
├── utils/            # Utilities
└── index.ts          # Main entry point
```

## License

ISC - Internal use only
