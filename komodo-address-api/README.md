# Komodo Address API

## Overview
Komodo Address API is a robust service designed to validate, normalize, and geocode postal addresses. It provides reliable address processing capabilities suitable for various applications that require accurate and standardized address data.

## Features
- **Validate**: Verify the correctness and existence of an address.
- **Normalize**: Standardize address formatting for consistency.
- **Geocode**: Convert addresses into geographic coordinates (latitude and longitude).

## Architecture Overview
The Komodo Address API is developed in Go and follows a modular folder structure:

- `internal/httpapi`: Contains the HTTP server and routing logic.
- `address`: Handles address validation and normalization logic.
- `geocode`: Responsible for geocoding services and integrations.
- `config`: Manages configuration and environment variables.

This structure facilitates maintainability and scalability of the service.

## Getting Started

### Prerequisites
- Go 1.22 or later
- Docker
- AWS CLI (optional, for deployment)

### Clone the repository
```bash
git clone https://github.com/komodoplatform/komodo-address-api.git
cd komodo-address-api
```

### Build the project
```bash
go build -o komodo-address-api ./cmd/komodo-address-api
```

### Run locally
```bash
./komodo-address-api
```

## Local Development

### Using `go run`
Run the API server directly with:
```bash
go run ./cmd/komodo-address-api
```

### Using Docker
Build the Docker image:
```bash
docker build -t komodo-address-api .
```
Run the container:
```bash
docker run -p 8080:8080 komodo-address-api
```

### Using Docker Compose
Start the service with:
```bash
docker-compose up
```

## API Endpoints

### POST /validate
Validate an address.
**Request:**
```json
{
  "address": "1600 Amphitheatre Parkway, Mountain View, CA"
}
```
**Response:**
```json
{
  "valid": true,
  "message": "Address is valid."
}
```

### POST /normalize
Normalize an address.
**Request:**
```json
{
  "address": "1600 amphitheatre pkwy, mountain view, ca"
}
```
**Response:**
```json
{
  "normalized_address": "1600 Amphitheatre Parkway, Mountain View, CA 94043, USA"
}
```

### POST /geocode
Geocode an address.
**Request:**
```json
{
  "address": "1600 Amphitheatre Parkway, Mountain View, CA"
}
```
**Response:**
```json
{
  "latitude": 37.4224764,
  "longitude": -122.0842499
}
```

## Configuration

The API can be configured using the following environment variables:

- `PORT`: Port number the server listens on (default: 8080).
- `GEOCODER`: Geocoding service to use (e.g., `googlemaps`).
- `GOOGLE_MAPS_API_KEY`: API key for Google Maps geocoding service.

## Deployment

### Build Docker Image
```bash
docker build -t komodo-address-api .
```

### Push to AWS ECR
```bash
aws ecr get-login-password --region <region> | docker login --username AWS --password-stdin <aws_account_id>.dkr.ecr.<region>.amazonaws.com
docker tag komodo-address-api:latest <aws_account_id>.dkr.ecr.<region>.amazonaws.com/komodo-address-api:latest
docker push <aws_account_id>.dkr.ecr.<region>.amazonaws.com/komodo-address-api:latest
```

### Run on AWS Fargate
Configure your ECS task to use the pushed image and deploy it on Fargate for scalable, serverless container hosting.

## Health Check

### GET /health
Check the health status of the API.
**Response:**
```json
{
  "status": "healthy"
}
```

## License
This project is licensed under the MIT License.
