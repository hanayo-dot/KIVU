# KIVU — Aquaculture & Lake Intelligence Platform

**KIVU** is a lake-wide aquaculture, water-quality, and spatial intelligence platform built for fish-cage farmers on **Lake Victoria**. It helps farmers monitor cage health in real time, receive AI-driven risk alerts, view farm and regional analytics, and evaluate zones for expansion using satellite and geospatial data.

> Backend banner: *KIVU Aquaculture & Lake Intelligence Backend*

---

## Features

### Level 1 — Cage Monitoring
- Submit and retrieve water-quality telemetry (dissolved oxygen, temperature, pH, turbidity)
- Automatic triggering of the **AI Risk Engine** on new readings
- Historical readings and observations (fish stress, algal growth, etc.)

### Level 2 — Farm Management
- Farmer registration & JWT authentication (phone + password)
- Farm dashboard with aggregated metrics, cage status, and active alerts
- Multi-cage management per farm

### Level 3 — Regional Lake Intelligence
- Lake zones modeled as PostGIS polygons (GeoJSON)
- Environmental metrics and spatial queries over Lake Victoria sectors (e.g. Homa Bay, Mbita, Rusinga)

### Level 4 — Expansion Signals
- Suitability scoring for new cage placement
- Nearest-zone evaluation via PostGIS distance queries
- Rationale-backed expansion recommendations

### Cross-cutting
- **AI Risk Engine** — rule-based anomaly detection (hypoxia, temperature spikes, deteriorating trends) with risk scores, confidence, and recommendations
- **Alerts & SMS** — automatic alert creation and optional SMS notifications
- **Copernicus integration** — background sync of satellite environmental data (demo/live modes)
- Interactive web frontend (Analytics, Alerts, Farm Reports, Maps, Settings)
- Full OpenAPI 3.0 specification

---

## Tech Stack

| Layer        | Technology                                      |
|--------------|-------------------------------------------------|
| Backend      | Go 1.22, Gin, sqlx, JWT (golang-jwt), bcrypt    |
| Database     | PostgreSQL 15 + PostGIS 3.3                     |
| Auth         | JWT access + refresh tokens                     |
| Frontend     | Vanilla HTML / CSS / JS (static, served by Gin) |
| Geospatial   | PostGIS geography (points & polygons)           |
| External     | Copernicus (satellite), SMS client              |
| Container    | Docker Compose                                  |
| API Docs     | OpenAPI 3.0 (`docs/openapi.yaml`)               |

---

## Project Structure

```
KIVU/
├── docker-compose.yml          # Postgres (PostGIS) + Backend services
├── backend/
│   ├── cmd/server/main.go      # Application entrypoint
│   ├── config/                 # Environment-based configuration
│   ├── database/
│   │   ├── migrations/         # PostGIS schema (0001_init_postgis_schema.sql)
│   │   └── seed/seed.sql       # Demo farmers, farms, cages, readings, zones
│   ├── integrations/           # Copernicus client, SMS client
│   ├── internal/
│   │   ├── handlers/           # HTTP handlers (auth, cages, farms, lake, AI, alerts…)
│   │   ├── middleware/         # CORS, logging, JWT auth
│   │   ├── models/             # Domain models
│   │   └── services/           # Auth, AI Risk Engine, Alerts, Expansion, Copernicus sync
│   ├── routes/routes.go        # Route registration + static frontend serving
│   ├── Dockerfile
│   ├── go.mod / go.sum
│   └── …
├── frontend/
│   ├── index.html / analytics.html / alerts.html / maps.html / reports.html / settings.html / login.html
│   ├── css/style.css
│   ├── js/api.js / app.js
│   └── images/
└── docs/
    └── openapi.yaml            # Full REST API specification
```

---

## Quick Start

### Prerequisites
- Docker & Docker Compose
- (Optional) Go 1.22+ for local backend development
- (Optional) Copernicus credentials for live satellite data

### 1. Clone / enter the project

```bash
cd KIVU
```

### 2. Start with Docker Compose

```bash
docker compose up --build
```

This starts:
- **PostgreSQL + PostGIS** on port `5432` (DB: `kivu`, user/password: `kivu` / `kivu_password`)
- **Backend** on port `8080`

Migrations and seed data are applied automatically on first start.

### 3. Access the application

| Resource              | URL                                      |
|-----------------------|------------------------------------------|
| Web Frontend          | http://localhost:8080                    |
| API base              | http://localhost:8080/api/v1             |
| OpenAPI specification | `docs/openapi.yaml`                      |

### Demo credentials
Seed data includes a sample farmer. Check `backend/database/seed/seed.sql` for exact phone/password values used in development (commonly `+254…` style numbers and simple demo passwords).

---

## Configuration

Environment variables (defaults shown):

| Variable                     | Default / Notes                                      |
|------------------------------|------------------------------------------------------|
| `SERVER_PORT`                | `8080`                                               |
| `ENVIRONMENT` / `APP_ENV`    | `development`                                        |
| `DB_HOST` / `DB_PORT`        | `postgres` / `5432` (in Docker)                      |
| `DB_USER` / `DB_PASSWORD`    | `kivu` / `kivu_password`                             |
| `DB_NAME`                    | `kivu`                                               |
| `JWT_SECRET`                 | Development secret (change in production!)           |
| `JWT_EXPIRATION_HOURS`       | `24` (or `JWT_ACCESS_EXP_MINUTES`)                   |
| `COPERNICUS_CLIENT_ID`       | Demo placeholder                                     |
| `COPERNICUS_CLIENT_SECRET`   | Demo placeholder                                     |
| `USE_COPERNICUS_LIVE`        | `false`                                              |

Create a `.env` file or export variables before running `docker compose up` if you need to override defaults.

---

## API Overview

Base path: `/api/v1`

Key endpoint groups (see `docs/openapi.yaml` for full details):

- **Authentication** — `POST /auth/register`, `POST /auth/login`, token refresh
- **Cage Monitoring** — `GET/POST /cages/{id}/readings`
- **Farm Management** — `GET /farms/{id}/dashboard`
- **Lake Intelligence** — `GET /lake/zones`
- **Expansion** — `GET /expansion/suitability`
- **AI / Alerts** — analysis results and alert management

All protected routes require a Bearer JWT.

---

## AI Risk Engine

The risk engine evaluates the latest readings for a cage and produces:

- **Risk score** (0–100+)
- **Risk level** (`low` / medium / high / critical)
- **Triggered factors** (e.g. critical hypoxia DO < 3.0 mg/L, deteriorating 3-reading DO trend, temperature anomalies)
- **Recommendation** and confidence score
- Automatic alert creation when thresholds are exceeded

It is invoked automatically when new telemetry is posted.

---

## Geospatial Capabilities

- Cages stored as `GEOGRAPHY(Point, 4326)`
- Lake zones as polygons
- Spatial indexes (GiST)
- Nearest-zone lookup via `ST_Distance` for expansion suitability

Seed data includes sample cages and zones around Homa Bay, Mbita Channel, and Rusinga areas of Lake Victoria.

---

## Development Notes

- Backend module: `github.com/hanayo-dot/KIVU/backend`
- Static frontend is served directly by the Gin server
- Background Copernicus sync loop starts with the server
- SMS client is pluggable (demo mode available)
- Tests exist for auth middleware, AI risk engine, Copernicus client, and auth service

### Local Go development (without Docker for the app)

```bash
# Start only the database
docker compose up postgres -d

# From backend/
go run ./cmd/server
```

Ensure `DB_HOST=localhost` (and matching credentials) when running outside Compose.

---

## License & Attribution

Built as a hackathon / demonstration platform for sustainable aquaculture decision support on Lake Victoria.

---

*KIVU — smarter cages, healthier lakes.*
