# Anomaly Detector

An HTTP API service that validates incoming requests against predefined API models, detecting anomalies in query parameters, headers, and body fields.

## Features

- Store API endpoint models with expected parameter types and requirements
- Validate incoming requests against stored models
- Detect anomalies including missing required fields and type mismatches
- Separate healthcheck server for monitoring

## Running Locally

### Prerequisites

- Go 1.25.0 or higher

### Setup

1. Install dependencies:
```bash
go mod download
```

2. Run the server:
```bash
go run main.go (Or press F5 for debug mode)
```

The service will start two servers:
- **Main API Server**: `http://localhost:8080`
- **Healthcheck Server**: `http://localhost:2802`

### Configuration

Configure via environment variables or `.env` file:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | Main API server port |
| `SERVER_HOST` | `localhost` | Main API server host |
| `HEALTHCHECK_PORT` | `2802` | Healthcheck server port |

## API Endpoints

### Healthcheck

Check if the service is running:

```bash
curl http://localhost:2802/health
```

**Response:**
```json
{"status":"healthy"}
```

### Store API Models

Store one or more API endpoint models for validation.

**Endpoint:** `POST /models`

**Example:**
```bash
curl -X POST http://localhost:8080/models \
  -H "Content-Type: application/json" \
  -d '[
    {
      "path": "/api/users",
      "method": "GET",
      "query_params": [
        {
          "name": "user_id",
          "types": ["Int", "UUID"],
          "required": true
        },
        {
          "name": "include_deleted",
          "types": ["Boolean"],
          "required": false
        }
      ],
      "headers": [
        {
          "name": "Authorization",
          "types": ["Auth-Token"],
          "required": true
        }
      ],
      "body": []
    }
  ]'
```

**Response:**
```json
{
  "message": "models stored successfully"
}
```

**Supported Parameter Types:**
- `String`
- `Int`
- `Boolean`
- `List`
- `Date`
- `Email`
- `UUID`
- `Auth-Token`

### Validate Request

Validate an incoming request against a stored model.

**Endpoint:** `POST /validate`

**Example (Valid Request):**
```bash
curl -X POST http://localhost:8080/validate \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/api/users",
    "method": "GET",
    "query_params": [
      {
        "name": "user_id",
        "value": "550e8400-e29b-41d4-a716-446655440000"
      },
      {
        "name": "include_deleted",
        "value": true
      }
    ],
    "headers": [
      {
        "name": "Authorization",
        "value": "Bearer abc123xyz"
      }
    ],
    "body": []
  }'
```

**Response (No Anomalies):**
```json
{
  "valid": true
}
```

**Example (Invalid Request - Missing Required Field):**
```bash
curl -X POST http://localhost:8080/validate \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/api/users",
    "method": "GET",
    "query_params": [
      {
        "name": "include_deleted",
        "value": true
      }
    ],
    "headers": [
      {
        "name": "Authorization",
        "value": "Bearer abc123xyz"
      }
    ],
    "body": []
  }'
```

**Response (With Anomalies):**
```json
{
  "valid": false,
  "anomalies": [
    {
      "field": "query_params",
      "parameter_name": "user_id",
      "reason": "required parameter \"user_id\" is missing"
    }
  ]
}
```

**Example (Type Mismatch):**
```bash
curl -X POST http://localhost:8080/validate \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/api/users",
    "method": "GET",
    "query_params": [
      {
        "name": "user_id",
        "value": "not-a-valid-uuid"
      }
    ],
    "headers": [
      {
        "name": "Authorization",
        "value": "Bearer abc123xyz"
      }
    ],
    "body": []
  }'
```

**Response:**
```json
{
  "valid": false,
  "anomalies": [
    {
      "field": "query_params",
      "parameter_name": "user_id",
      "reason": "type mismatch: expected one of [Int UUID] types, but got the type string"
    }
  ]
}
```

## Architecture

- **Dependency Injection**: Uses `uber/dig` for IoC container
- **Routing**: Gorilla Mux for HTTP routing
- **Graceful Shutdown**: Handles SIGINT/SIGTERM signals
- **Logging**: Structured JSON logging with `slog`

## Design Tradeoffs

### 1. No Dockerfile

The project does not include a Dockerfile for containerization.

**Tradeoff:**

**Advantages:**
- **Simpler local development**: Just `go run main.go` - no Docker installation needed
- **Faster iteration**: No container build/rebuild steps during development
- **Lower barrier to entry**: Only Go required, not Docker daemon/CLI

**Disadvantages:**
- **No deployment standardization**: "Works on my machine" problems across environments
- **Manual environment setup**: Developers need correct Go version and dependencies
- **Harder cloud deployment**: Most platforms expect containers (Kubernetes, ECS - AWS, Cloud Run - GCP)

**Why this matters:**
This is a **development simplicity vs production readiness** tradeoff. For local development and prototyping, Docker adds overhead that is not always necessary. For production deployments, a Dockerfile becomes essential for reproducible builds and deployment to container orchestration platforms.

### 2. All-or-Nothing Model Storage

When storing multiple API models via `POST /models`, the implementation validates ALL models first before storing ANY of them. This is an intentional design decision with the following tradeoffs:

**Approach:**
1. First loop: Validate all models (check for nil, empty fields, duplicates)
2. Acquire lock
3. Second loop: Store all models if validation passed

**Advantages:**
- **Atomicity**: Either all models are stored successfully or none are stored, preventing partial updates
- **Consistency**: The store never enters an inconsistent state where some models from a batch are stored and others are rejected

**Disadvantages:**
- **Double iteration**: Models are processed twice (validation + storage), though this is negligible for typical batch sizes

**Why this matters:**
If a request contains 100 models and the 50th is invalid, an alternative approach might store the first 49 before failing. This would require either manual rollback or leave the system in a partially updated state. The all-or-nothing approach treats each request as a transaction, which is more predictable for API consumers.

### 3. Separate Healthcheck Server

The application runs two HTTP servers: a main API server (port 8080) and a separate healthcheck server (port 2802), instead of combining them into a single server.

**Tradeoff:**

**Advantages:**
- **Isolation**: If the main server hangs (deadlock, slow requests, resource exhaustion), the healthcheck can still respond, allowing orchestrators (Kubernetes, Docker) to detect the issue
- **Security**: Healthcheck bypasses main server routing/middleware/auth

**Disadvantages:**
- **More resources**: Two servers means two listening sockets and additional goroutines
- **Configuration overhead**: Managing two ports instead of one

**Why this matters:**
Orchestration platforms use health endpoints to detect when a service is degraded vs completely down. With a combined server, if the main API is blocked handling slow requests, health checks fail too, causing unnecessary restarts. Separate servers ensure health monitoring works independently.

### 4. In-Memory Storage

API models are currently stored in-memory using a Go map, rather than in a database.

**Tradeoff:**

**Advantages:**
- **Performance**: Extremely fast reads/writes (nanoseconds vs milliseconds)
- **Simplicity**: No database setup, connection pools, or migration scripts
- **Zero dependencies**: No external services required

**Disadvantages:**
- **No persistence**: Data lost on restart
- **Limited by RAM**: Cannot store unlimited models
- **No horizontal scaling**: Each instance has separate data

**Future migration path:**

If persistence becomes necessary, **PostgreSQL** can be a good choice because:
- Native JSON storage with indexing on `(path, method)` for fast lookups
- ACID transactions align with the all-or-nothing storage approach
- Room for growth: model versioning, soft deletes, audit logs
- The `IModelStore` interface allows adding a `PostgresModelStore` implementation without changing consumers

### 5. Parameter Validation with HashMap Lookup

During validation, request parameters are first converted into a hash map for O(1) lookups instead of repeatedly iterating through the parameter list.

**Tradeoff:**
- **Time**: O(n + m) instead of O(n × m) where n = model params, m = request params
- **Space**: Additional O(m) memory for the hash map
- **Example**: 20 model params × 50 request params = up to 1,000 iterations vs 70 with hash map

**Why this matters:**
Requests with many query parameters and headers benefit from constant-time lookups. The memory overhead is minimal and short-lived (deallocated after validation).