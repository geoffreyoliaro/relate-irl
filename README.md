# MiniNexus — Relationship Intelligence API

> A demo platform built for the Orbiter.io interview.
> Stack: **Go + Gin** · **FalkorDB** (graph DB) · **Xano** (auth) · **GCP Cloud Run** · **GitHub Actions**

---

## What this demonstrates

| Skill | How it shows up |
|---|---|
| Graph DB design | FalkorDB Cypher queries for shortest path, mutual connections, relationship scoring |
| Go API design | Gin, clean handler/model/graph layering, context propagation |
| Auth architecture | Xano validates credentials → Go issues short-lived JWT (decoupled) |
| GCP deployment | Cloud Run (serverless containers), Artifact Registry, Workload Identity Federation |
| CI/CD | GitHub Actions: test → build → push → deploy → smoke test |
| Local dev | Docker Compose with FalkorDB + hot-reload |

---

## Quick start (local)

```bash
# 1. Start FalkorDB
docker compose up falkordb -d

# 2. Copy env and run
cp .env.example .env
go run ./cmd/api

# Or run everything with Docker Compose
docker compose up --build
```

The API starts on **http://localhost:8080**.  
Demo data is seeded automatically on startup.

---

## API reference

### Auth
```
POST /auth/login
{ "email": "amara@horizonvc.com", "password": "demo123" }
→ { "token": "eyJ...", "user": {...} }
```

### People
```
GET    /api/v1/people
POST   /api/v1/people
GET    /api/v1/people/:id          # person + their 1-hop network
```

### Relationships
```
POST   /api/v1/relationships
{ "from_id": "p1", "to_id": "p3", "type": "advisor", "strength": 70 }
```

### Intelligence queries (demo showpieces)
```
# Who do both Amara (p1) and Sofia (p4) know?
GET /api/v1/intelligence/mutual-connections?a=p1&b=p4

# Shortest connection chain from Amara to Chioma
GET /api/v1/intelligence/shortest-path?from=p1&to=p8

# How strong is the relationship between James and Sofia?
GET /api/v1/intelligence/relationship-strength?a=p3&b=p4
```

---

## FalkorDB key concepts (for the interview)

**Why FalkorDB over Postgres/Neo4j?**
- Built on Redis — sub-millisecond latency for graph traversals
- Native Cypher support (same query language as Neo4j)
- Sparse, evolving schemas suit relationship data naturally
- Relationship properties (strength, type, timestamp) are first-class

**Key Cypher patterns used:**

```cypher
-- Mutual connections (2-hop traversal)
MATCH (a:Person {id: $idA})-[:KNOWS]-(mutual:Person)-[:KNOWS]-(b:Person {id: $idB})
WHERE a <> b
RETURN mutual.name, mutual.company

-- Shortest path (up to 6 hops)
MATCH p = shortestPath((a:Person {id: $from})-[:KNOWS*..6]-(b:Person {id: $to}))
RETURN [node IN nodes(p) | node.name] AS path, length(p) AS hops

-- Composite relationship strength
WITH COALESCE(direct.strength, 0) AS direct_strength,
     COUNT(DISTINCT mutual) AS mutual_count
RETURN (direct_strength * 0.6 + mutual_count * 10) AS composite_score
```

---

## GCP deployment

### One-time setup
```bash
# Create Artifact Registry repo
gcloud artifacts repositories create mininexus \
  --repository-format=docker \
  --location=us-central1

# Set up Workload Identity Federation (no service account keys needed)
# See: https://github.com/google-github-actions/auth#setup
```

### GitHub Secrets required
```
GCP_PROJECT_ID
GCP_WORKLOAD_IDENTITY_PROVIDER
GCP_SERVICE_ACCOUNT
FALKORDB_HOST
FALKORDB_PASSWORD
JWT_SECRET
XANO_API_BASE
```

### Manual deploy (without Actions)
```bash
gcloud run deploy mininexus-api \
  --image us-central1-docker.pkg.dev/YOUR_PROJECT/mininexus/api:latest \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars FALKORDB_HOST=...,JWT_SECRET=...
```

---

## Architecture decisions to talk through

1. **Xano for auth, Go for graph logic** — Xano handles user management, email verification, and password reset without custom code. The Go service is purely graph/intelligence logic. This mirrors how real startups compose tooling.

2. **JWT decoupling** — After Xano validates, we issue our own short-lived JWT. This means the graph API has no runtime dependency on Xano for every request — just at login time.

3. **Cloud Run over GKE** — No persistent state in the Go service (FalkorDB is external), so scale-to-zero is safe and dramatically cuts cost. GKE adds overhead without benefit here.

4. **Workload Identity Federation** — No long-lived service account JSON keys in GitHub Secrets. GCP issues short-lived tokens per workflow run.

5. **Relationship strength formula** — `composite = (direct_strength × 0.6) + (mutual_count × 10)`. Direct edges weighted more than shared contacts, but shared contacts provide signal when there's no direct edge. Easy to tune.

---

## Demo seed data

| ID | Name | Company | Role |
|---|---|---|---|
| p1 | Amara Osei | Horizon VC | Partner |
| p2 | Lena Fischer | BuildCo | CEO |
| p3 | James Mwangi | TechBridge | CTO |
| p4 | Sofia Reyes | OpenLoop AI | Founder |
| p5 | Tom Nakamura | Horizon VC | Analyst |
| p6 | Priya Sharma | BuildCo | VP Engineering |
| p7 | David Oliaro | NexusTech | Head of Growth |
| p8 | Chioma Eze | OpenLoop AI | COO |

Good demo queries:
- `mutual?a=p1&b=p4` → who do the VC and the founder both know?
- `shortest-path?from=p1&to=p8` → 3-hop chain: Amara → Sofia → Chioma
- `strength?a=p2&b=p4` → no direct edge, but 2 mutual contacts
