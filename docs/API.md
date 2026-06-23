# Relate IRL - API Documentation

## Base URL

Development: `http://localhost:8080`
Production: `https://relate-irl-api-xxxxx.run.app`

## Authentication

All endpoints except `/auth/login` require a JWT token in the `Authorization` header:

```
Authorization: Bearer <jwt_token>
```

## Endpoints

### Public Endpoints

#### POST /auth/login
Authenticate with Xano and receive a JWT token.

**Request:**
```json
{
  "email": "amara@horizonvc.com",
  "password": "demo123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "user": {
    "id": 123,
    "name": "Amara",
    "email": "amara@horizonvc.com"
  }
}
```

**Error (401 Unauthorized):**
```json
{
  "error": "invalid credentials"
}
```

#### GET /healthz
Health check endpoint for Cloud Run and load balancers.

**Response (200 OK):**
```json
{
  "status": "ok"
}
```

---

### Protected Endpoints

All below endpoints require JWT authentication.

#### GET /api/v1/people
List all people in the network.

**Response (200 OK):**
```json
{
  "people": [
    {
      "id": "person:1",
      "name": "Alice Johnson",
      "title": "CEO",
      "company": "TechCorp",
      "email": "alice@techcorp.com"
    }
  ]
}
```

#### POST /api/v1/people
Create a new person.

**Request:**
```json
{
  "name": "Bob Smith",
  "title": "Product Manager",
  "company": "StartupXYZ",
  "email": "bob@startupxyz.com"
}
```

**Response (201 Created):**
```json
{
  "id": "person:42",
  "name": "Bob Smith",
  "title": "Product Manager",
  "company": "StartupXYZ",
  "email": "bob@startupxyz.com"
}
```

#### GET /api/v1/people/:id
Get details of a specific person.

**Response (200 OK):**
```json
{
  "id": "person:1",
  "name": "Alice Johnson",
  "title": "CEO",
  "company": "TechCorp",
  "email": "alice@techcorp.com"
}
```

#### POST /api/v1/relationships
Create a relationship between two people.

**Request:**
```json
{
  "from_id": "person:1",
  "to_id": "person:2",
  "type": "colleague",
  "strength": 0.8,
  "notes": "Worked together on Project X"
}
```

**Response (201 Created):**
```json
{
  "id": "relationship:42",
  "from_id": "person:1",
  "to_id": "person:2",
  "type": "colleague",
  "strength": 0.8,
  "notes": "Worked together on Project X"
}
```

#### GET /api/v1/relationships/:id/network
Get the network around a relationship node.

**Query Parameters:**
- `depth` (optional, default=2): How many hops to traverse
- `limit` (optional, default=100): Max nodes to return

**Response (200 OK):**
```json
{
  "center": {
    "id": "person:1",
    "name": "Alice Johnson",
    "title": "CEO"
  },
  "connections": [
    {
      "id": "person:2",
      "name": "Bob Smith",
      "distance": 1,
      "via": "colleague"
    }
  ]
}
```

#### GET /api/v1/intelligence/mutual-connections
Find mutual connections between two people.

**Query Parameters:**
- `person_a`: First person ID (required)
- `person_b`: Second person ID (required)

**Response (200 OK):**
```json
{
  "person_a": "person:1",
  "person_b": "person:2",
  "mutual_connections": [
    {
      "id": "person:3",
      "name": "Charlie Brown",
      "connection_path": ["person:1", "person:3", "person:2"]
    }
  ],
  "count": 1
}
```

#### GET /api/v1/intelligence/shortest-path
Find the shortest path between two people.

**Query Parameters:**
- `from`: Starting person ID (required)
- `to`: Target person ID (required)

**Response (200 OK):**
```json
{
  "from": "person:1",
  "to": "person:5",
  "path": [
    {"id": "person:1", "name": "Alice"},
    {"id": "person:2", "name": "Bob"},
    {"id": "person:5", "name": "Eve"}
  ],
  "length": 3,
  "strength": 0.72
}
```

#### GET /api/v1/intelligence/relationship-strength
Calculate overall relationship strength.

**Query Parameters:**
- `from`: Person A ID (required)
- `to`: Person B ID (required)

**Response (200 OK):**
```json
{
  "from": "person:1",
  "to": "person:2",
  "strength": 0.85,
  "factors": {
    "direct_connection": 0.9,
    "path_quality": 0.8,
    "mutual_connections": 1.0
  }
}
```

---

## Error Responses

### 400 Bad Request
Invalid request parameters.

```json
{
  "error": "invalid request"
}
```

### 401 Unauthorized
Missing or invalid JWT token.

```json
{
  "error": "unauthorized"
}
```

### 404 Not Found
Resource not found.

```json
{
  "error": "not found"
}
```

### 500 Internal Server Error
Server error.

```json
{
  "error": "internal error"
}
```

---

## Example Flow

### 1. Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"amara@horizonvc.com","password":"demo123"}'
```

### 2. Use token for protected endpoints
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl http://localhost:8080/api/v1/people \
  -H "Authorization: Bearer $TOKEN"
```

### 3. Query intelligence endpoints
```bash
curl "http://localhost:8080/api/v1/intelligence/mutual-connections?person_a=person:1&person_b=person:2" \
  -H "Authorization: Bearer $TOKEN"
```

---

## Testing with curl

```bash
# Health check
curl http://localhost:8080/healthz

# Login (get token)
RESPONSE=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"amara@horizonvc.com","password":"demo123"}')

TOKEN=$(echo $RESPONSE | jq -r '.token')

# List people (requires token)
curl http://localhost:8080/api/v1/people \
  -H "Authorization: Bearer $TOKEN"
```

---

## Rate Limiting

Currently not implemented. Will be added in future releases.

## CORS

The API is configured to accept requests from any origin. In production, this should be restricted.

## Versioning

API version is included in the URL: `/api/v1/`

Future versions will use `/api/v2/`, etc.
