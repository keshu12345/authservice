# Authentication Service API Documentation

## Overview

This document provides detailed information about the authentication service API endpoints, including request formats, required headers, and expected responses.

## Getting Started

### Environment Configurations

The application supports three different environments:

1. **Local Development**: `config/local`
2. **Non-Production**: `config/nonprod`
3. **Production**: `config/prod`

### Running the Application

#### Option 1: Run directly with Go

To start the authentication service, use the following command with the appropriate configuration:

```bash
go run cmd/main.go -config config/local -migrations migrations
```

This command:
- Runs the main application entry point
- Loads configuration from the specified environment directory (e.g., `config/local`)
- Applies database migrations from the `migrations` directory

To use a different environment configuration:
```bash
go run cmd/main.go -config config/nonprod -migrations migrations
```

or

```bash
go run cmd/main.go -config config/prod -migrations migrations
```

#### Option 2: Run with Docker Compose

For a complete environment including the database, pgAdmin, and the auth service:

```bash
docker-compose up
```

This will start all services defined in the Docker Compose file, and the auth service will be available at http://localhost:8080.

### Docker Setup

#### Docker Compose

The service uses Docker Compose for local development. Create a `docker-compose.yml` file with the following content:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    container_name: my_postgres
    restart: always
    environment:
      POSTGRES_USER: myuser
      POSTGRES_PASSWORD: mypassword
      POSTGRES_DB: mydatabase
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    command: postgres -c max_connections=100

  pgadmin:
    image: dpage/pgadmin4:7
    container_name: my_pgadmin
    restart: always
    environment:
      PGADMIN_DEFAULT_EMAIL: admin@local.com
      PGADMIN_DEFAULT_PASSWORD: admin123
    ports:
      - "8081:80"
    depends_on:
      - postgres
    volumes:
      - pgadmin_data:/var/lib/pgadmin

  authservice:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: my_authservice
    depends_on:
      - postgres
    ports:
      - "8080:8080"
    environment:
      GIN_MODE: release
    volumes:
      - ./logs:/app/logs

volumes:
  pgdata:
  pgadmin_data:
```

This will:
- Start a PostgreSQL 15 database on port 5432
- Start pgAdmin on port 8081 (accessible at http://localhost:8081)
- Build and start the auth service on port 8080
- Configure database credentials (User: myuser, Password: mypassword, Database: mydatabase)
- Configure pgAdmin credentials (Email: admin@local.com, Password: admin123)

#### Dockerfile

The service uses a multi-stage build process for optimal container size:

```dockerfile
# Use Go 1.24.2 to satisfy go.mod requirement
FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o authservice ./cmd/main.go

# --- Runtime Stage ---
FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata netcat-openbsd

COPY --from=builder /app/authservice .
COPY start.sh .
COPY config ./config
COPY migrations ./migrations

# Fix line endings and permissions
RUN sed -i 's/\r$//' start.sh && chmod +x start.sh authservice

ENTRYPOINT ["/app/start.sh"]
```

### Base URL

Once running, the service is available at:

```
http://localhost:8080
```

## Authentication Methods

The API uses two authentication methods:

1. **API Key Authentication**: Required for all API endpoints
   - Send via `X-API-Key` header

2. **JWT Authentication**: Required for protected routes
   - Send via `Authorization` header with format: `Bearer <token>`

## API Endpoints

### Authentication

#### Register User (Sign Up)

Creates a new user account with specified role.

- **URL**: `/v1/signup`
- **Method**: `POST`
- **Auth Required**: None
- **Content-Type**: `application/json`

**Request Body**:
```json
{
  "username": "keshav",
  "password": "password123",
  "role": "creator"
}

```

**Success Response**:
- **Code**: 201 Created
- **Content**:
```json
{
    "timestamp": "2025-05-03T19:05:34+05:30",
    "code": 201,
    "status": "Created",
    "message": "User registered successfully",
    "data": {
        "username ": "keshav"
    }
}
```

**Error Responses**:
- **Code**: 400 Bad Request
  - User already exists
  - Invalid role
  - Missing required fields
- **Code**: 500 Internal Server Error

**Example**:
```bash
curl -X POST -H "Content-Type: application/json" \
     -d '{"username":"alice_creator","password":"password123","role":"creator"}' \
     http://localhost:8080/v1/signup
```

#### Login User (Sign In)

Authenticates a user and returns a JWT token.

- **URL**: `/v1/signin`
- **Method**: `POST`
- **Auth Required**: None
- **Content-Type**: `application/json`

**Request Body**:
```json
  {
    "username": "keshav",
    "password": "password123"
  }

```

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
  {
      "timestamp": "2025-05-03T19:10:50+05:30",
      "code": 200,
      "status": "OK",
      "message": "Login successful",
      "data": {
          "token ": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJ1c2VyLWNyZWF0b3ItMjAyNTA1MDMxNzI5MDciLCJyb2xlcyI6WyJjcmVhdG9yIl0sImlzcyI6ImF1dGhzZXJ2aWNlIiwic3ViIjoia2VzaGF2IiwiZXhwIjoxNzQ2MjgzMjUwLCJuYmYiOjE3NDYyNzk2NTAsImlhdCI6MTc0NjI3OTY1MCwianRpIjoiMTgzYzA4MTcxZWQwZjZkOCJ9.CBUQo8n4EMEWIDq7yflahwxspWwPCmi9UY7sYfffTDU",
          "apiKey": "8a3575e58a29cd675ed310569e6def9d1b5e81e1ec647bc05eb263b057fcd5c6",
          "userId ": "user-creator-20250503172907",
          "role ": "creator"
      }
  }
```

**Error Responses**:
- **Code**: 401 Unauthorized
  - Invalid credentials
- **Code**: 500 Internal Server Error

**CURL**:
```bash
curl -X POST -H "Content-Type: application/json" \
     -d '{"username":"alice_creator","password":"password123"}' \
     http://localhost:8080/v1/signin
```

#### Refresh Token

Generates a new JWT token using an existing valid token.

- **URL**: `/v1/refresh`
- **Method**: `POST`
- **Auth Required**: None
- **Content-Type**: `application/json`

**Request Body**:
```json
{
  "token": "string" // Valid JWT token
}
```

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "timestamp": "2025-05-03T10:38:15Z",
  "code": 200,
  "status": "OK",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJ1c2VyLWNyZWF0b3ItMjAyNTA1MDMxMDM0MTUiLCJyb2xlcyI6WyJjcmVhdG9yIl0sImlzcyI6ImF1dGhzZXJ2aWNlIiwic3ViIjoiYWxpY2VfY3JlYXRvciIsImV4cCI6MTc0NjI3MjI5NSwibmJmIjoxNzQ2MjY4Njk1LCJpYXQiOjE3NDYyNjg2OTUsImp0aSI6IjE4M2JmZTIwNzE0MzU1OGEifQ.B880WOZh5P3-Ok2f8GKvEDENANFy_bAbDINdLNGJ6t4",
    "refreshToken": "",
    "apiKey": "16ddddcbe532e93179681b8ba0d4ff7b73d211b6c8a30d2772b48bf9356e92b9",
    "userId": "user-creator-20250503103415",
    "role": "creator"
  }
}
```

**Error Responses**:
- **Code**: 401 Unauthorized
  - Invalid or expired token
- **Code**: 500 Internal Server Error

**CURL**:
```bash
curl -X POST http://localhost:8080/v1/refresh \
  -H "Content-Type: application/json" \
  -d '{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}'
```

#### Revoke Token

Revokes a JWT token, making it invalid for future requests.

- **URL**: `/v1/revoke`
- **Method**: `POST`
- **Auth Required**: JWT
- **Headers**:
  - `Authorization: Bearer <token>`
- **Content-Type**: `application/json`

**Success Response**:
- **Code**: 202 Accepted
- **Content**:
```json
  {
    "timestamp": "2025-05-03T10:43:49Z",
    "code": 202,
    "status": "Accepted",
    "message": "Token revoked successfully"
  }
```

**Error Responses**:
- **Code**: 401 Unauthorized
  - Invalid token
- **Code**: 500 Internal Server Error

**Example**:
```bash
curl -X POST http://localhost:8080/v1/revoke \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

#### Unrevoke Token

Unrevoking a previously revoked JWT token, making it valid again.

- **URL**: `/v1/unrevoke`
- **Method**: `POST`
- **Auth Required**: JWT
- **Headers**:
  - `Authorization: Bearer <token>`
- **Content-Type**: `application/json`

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "timestamp": "2025-05-03T10:46:33Z",
  "code": 200,
  "status": "OK",
  "message": "Token unrevoked successfully"
}
```

**Error Responses**:
- **Code**: 401 Unauthorized
  - Invalid token
- **Code**: 500 Internal Server Error

**CURL**:
```bash
curl -X POST http://localhost:8080/v1/unrevoke \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

### Protected Endpoints

#### Get Creator Information

Retrieves information about creators. Requires both API key and valid JWT token with creator role.

- **URL**: `/api/creators`
- **Method**: `GET`
- **Auth Required**: API Key + JWT
- **Headers**:
  - `X-API-Key: <api_key>`
  - `Authorization: Bearer <token>`

**Success Response**:
- **Code**: 200 OK
- **Content**:
```json
{
  "timestamp": "2025-05-03T10:49:54Z",
  "code": 200,
  "status": "OK",
  "message": "User found",
  "data": {
    "userId": "user-creator-20250503103415",
    "username": "keshav",
    "password_hash": "$2a$10$XSVv8fGxib5AB9Q4hUNe.etywfdisU2nLinLVHxZluEcpC1yR944S",
    "role": "creator",
    "apiKey": "16ddddcbe532e93179681b8ba0d4ff7b73d211b6c8a30d2772b48bf9356e92b9",
    "is_revoke": false
  }
}
```

**Error Responses**:
- **Code**: 401 Unauthorized
  - Invalid API Key
  - Invalid JWT Token
- **Code**: 403 Forbidden
  - Insufficient permissions (missing creator role)
- **Code**: 500 Internal Server Error

**Example**:
```bash
curl -H "X-API-Key: 26355af1f816cd4bffaf1b62c040369a92c523c2055a8b1a0be44b47e2fdb426" \
     -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
     http://localhost:8080/api/creators
```

## API Keys

The API uses predefined API keys for authentication: these api keys are generated at time create user 
 
| Role    | API Key                                                           |
|---------|-------------------------------------------------------------------|
| Creator | 26355af1f816cd4bffaf1b62c040369a92c523c2055a8b1a0be44b47e2fdb426 |
| Partner | 89831df8cd8186406c677eaabb83df6cd2fe7a38b80dd170329d0c086d5296d3 |

## Error Handling

The API returns appropriate HTTP status codes along with error messages in the response body:

```json
{
  "timestamp": "2025-05-03T10:45:12Z",
  "code": 401,
  "status": "Unauthorized",
  "message": "Error message description",
  "data": null
}
```

Common error codes:
- 400: Bad Request - Invalid input parameters
- 401: Unauthorized - Authentication failed
- 403: Forbidden - Insufficient permissions
- 404: Not Found - Resource not found
- 500: Internal Server Error - Server-side issue

## Middleware

The API uses the following middleware components:

1. **LoggingMiddleware**: Logs all incoming requests
2. **ApiKeyAuthMiddleware**: Validates API key for all API endpoints
3. **JwtAuthMiddleware**: Validates JWT token for protected routes