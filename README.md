# Safiri Logistics

Safiri Logistics is a backend platform for heavy-goods freight movement: machinery, vehicles, timber, industrial equipment, and other trailer-based cargo. It supports authentication, driver KYC, distance-based order matching, admin dispatch optimization, real-time tracking, freight pricing, and M-Pesa payment initiation.

## Why This System Looks The Way It Does

This project is intentionally designed as a modular monolith first, not a microservice fleet.

The design rationale is:

- keep latency low for the first matching and dispatch workflows
- keep deployment simple while the business rules are still evolving
- make assignments explainable instead of hiding dispatch logic in a black box
- persist critical business events early: loads, assignments, tracking, payments
- introduce clean seams for Redis, Kafka, and AI later without over-engineering now

In practical terms, that means:

- PostgreSQL is the source of truth
- Gin handles the HTTP API and WebSocket upgrade path
- matching, pricing, payment orchestration, and tracking live in separate internal packages
- dispatch logic is heuristic and auditable
- real-time updates are supported today, while Redis/Kafka remain future scale levers

If you want the longer design review, see [system_design.md](/c:/Users/HP%20ENVY/Desktop/Projects/safiri-logistics/docs/system_design.md).

## Core Capabilities

- JWT authentication for customers, drivers, and admins
- Driver KYC submission and admin approval
- Driver operational state:
  - online or offline
  - available or busy
  - location
  - experience
  - equipment type
  - max supported load
- Heavy-goods load posting with:
  - pickup and dropoff coordinates
  - pickup schedule
  - priority
  - equipment requirement
  - cargo type
- Automatic distance-based driver matching
- Admin batch dispatch optimization for future loads
- Real-time tracking via WebSockets
- Freight pricing engine
- M-Pesa checkout initiation

## High-Level Architecture

```mermaid
flowchart LR
    C[Customer App] --> API[Safiri API Service]
    D[Driver App] --> API
    A[Admin Ops Console] --> API

    API --> AUTH[Auth Module]
    API --> DRV[Driver Ops Module]
    API --> ORD[Order Module]
    API --> DSP[Dispatch Engine]
    API --> TRK[Tracking Hub]
    API --> PRC[Pricing Engine]
    API --> PAY[Payment Orchestrator]

    AUTH --> DB[(PostgreSQL)]
    DRV --> DB
    ORD --> DB
    DSP --> DB
    TRK --> DB
    PAY --> DB

    TRK --> WS[WebSocket Clients]
    PAY --> MPESA[M-Pesa Provider Adapter]

    DSP -. future .-> REDIS[(Redis)]
    API -. future .-> KAFKA[(Kafka)]
```

### What Each Piece Does

- `Auth Module`: registration, login, JWT issuance, user context
- `Driver Ops Module`: KYC, operational state, live driver readiness
- `Order Module`: load creation, status transitions, ownership checks
- `Dispatch Engine`: nearest eligible driver matching and admin optimization
- `Tracking Hub`: persists location pings and broadcasts live updates
- `Pricing Engine`: deterministic freight quote calculation
- `Payment Orchestrator`: initiates provider checkout and stores payment intent

## Matching Philosophy

Matching is optimized for explainability and low latency.

The current engine:

1. filters drivers by KYC, availability, equipment, capacity, and live location
2. scores candidates by distance, experience, and load priority
3. assigns the best candidate
4. stores the assignment score for auditability

The scoring model is simple on purpose:

`score = distance_km * 0.65 + experience_penalty * 0.20 + priority_penalty * 0.15`

Why this is a good first step:

- fast enough for automatic matching on load creation
- understandable to operations teams
- easy to test
- easy to improve later with Redis geo lookup or optimization solvers

## Database Relationships

```mermaid
erDiagram
    USERS ||--o{ LOADS : posts
    USERS ||--|| DRIVER_PROFILES : owns
    USERS ||--o{ TRACKING_EVENTS : sends
    USERS ||--o{ PAYMENTS : initiates
    USERS ||--o{ LOADS : assigned_driver

    LOADS ||--o{ TRACKING_EVENTS : has
    LOADS ||--o{ PAYMENTS : has

    USERS {
        uuid id PK
        text name
        text email
        text password_hash
        text role
        timestamp created_at
    }

    DRIVER_PROFILES {
        uuid user_id PK, FK
        text national_id
        text license_number
        text truck_registration
        text status
        int years_experience
        numeric max_load_kg
        text equipment_type
        bool is_online
        bool is_available
        double current_latitude
        double current_longitude
        timestamp last_location_at
    }

    LOADS {
        uuid id PK
        uuid poster_id FK
        uuid assigned_driver_id FK
        text title
        text description
        text origin
        text destination
        numeric weight_kg
        text status
        double pickup_latitude
        double pickup_longitude
        double dropoff_latitude
        double dropoff_longitude
        timestamp pickup_at
        int priority
        text cargo_type
        text equipment_type
        numeric quoted_price_kes
        double matching_score
        text assignment_source
    }

    TRACKING_EVENTS {
        uuid id PK
        uuid load_id FK
        uuid driver_id FK
        double latitude
        double longitude
        double speed_kph
        double heading_degrees
        timestamp recorded_at
    }

    PAYMENTS {
        uuid id PK
        uuid load_id FK
        uuid customer_id FK
        text provider
        text provider_reference
        text phone_number
        numeric amount_kes
        text currency
        text status
        timestamp created_at
        timestamp updated_at
    }
```

### Key Relationship Notes

- A `user` can be a customer, driver, or admin.
- A driver has exactly one `driver_profile`.
- A customer can post many `loads`.
- A load can be assigned to one driver at a time.
- A load can have many `tracking_events`.
- A load can have one or more `payments` over time, depending on how billing evolves.

## API Overview

Base URL:

```text
http://localhost:8080/api/v1
```

Authentication:

- Protected routes require `Authorization: Bearer <token>`

### Auth

`POST /auth/register`

- Registers a `customer` or `driver`

Example body:

```json
{
  "name": "Customer One",
  "email": "customer@example.com",
  "password": "Password123",
  "role": "customer"
}
```

`POST /auth/login`

Example body:

```json
{
  "email": "customer@example.com",
  "password": "Password123"
}
```

`GET /auth/me`

- Returns the currently authenticated user

### Driver

`POST /drivers/kyc`

- Driver submits KYC details

Example body:

```json
{
  "national_id": "12345678",
  "license_number": "DL-90876",
  "truck_registration": "KDA123A"
}
```

`GET /drivers/kyc/me`

- Returns the authenticated driver's KYC profile

`PATCH /drivers/:userID/kyc/review`

- Admin approves or rejects driver KYC

Example body:

```json
{
  "status": "approved"
}
```

`PATCH /drivers/me/operations`

- Driver updates availability, location, capacity, equipment, and experience

Example body:

```json
{
  "years_experience": 8,
  "max_load_kg": 34000,
  "equipment_type": "lowbed",
  "is_online": true,
  "is_available": true,
  "latitude": -1.286389,
  "longitude": 36.817223
}
```

### Loads

`POST /loads`

- Customer creates a load
- System calculates `quoted_price_kes`
- System attempts automatic distance-based matching

Example body:

```json
{
  "title": "Excavator Move",
  "description": "Lowbed movement for industrial equipment",
  "origin": "Mombasa",
  "destination": "Nairobi",
  "weight_kg": 12000,
  "pickup_latitude": -4.043477,
  "pickup_longitude": 39.668206,
  "dropoff_latitude": -1.286389,
  "dropoff_longitude": 36.817223,
  "pickup_at": "2026-04-11T08:00:00Z",
  "priority": 5,
  "cargo_type": "machinery",
  "equipment_type": "lowbed"
}
```

`GET /loads`

- Customer sees own loads
- Driver or admin sees broader load list

`POST /loads/:loadID/pick`

- Assigned or eligible driver confirms pickup

`PATCH /loads/:loadID/status`

- Driver moves load through shipment lifecycle

Example body:

```json
{
  "status": "in_transit"
}
```

Valid transitions today:

- `posted -> picked`
- `matched -> picked`
- `picked -> in_transit`
- `in_transit -> delivered`

### Dispatch

`POST /dispatch/optimize`

- Admin matches future loads in bulk for a given pickup day

Example body:

```json
{
  "pickup_date": "2026-04-11T00:00:00Z"
}
```

### Tracking

`POST /tracking/loads/:loadID/events`

- Assigned driver posts a tracking ping

Example body:

```json
{
  "latitude": -2.154007,
  "longitude": 37.990985,
  "speed_kph": 62,
  "heading_degrees": 135
}
```

`GET /tracking/ws?load_id=:loadID`

- WebSocket stream for live tracking updates

### Payments

`POST /payments/loads/:loadID/mpesa-checkout`

- Customer initiates M-Pesa checkout for a quoted load

Example body:

```json
{
  "phone_number": "254700000001"
}
```

## Status Model

Current load statuses:

- `posted`
- `matched`
- `picked`
- `in_transit`
- `delivered`

Current payment statuses:

- `initiated`
- `paid`
- `failed`

Current KYC statuses:

- `pending`
- `approved`
- `rejected`

## Local Development

### Environment

`.env` should contain:

```env
DATABASE_URL=postgresql://postgres:password@localhost:5432/safiri?sslmode=disable
JWT_SECRET=change-me
PORT=8080
JWT_TOKEN_TTL_HOURS=24
```

### Run Migrations

Apply these in order:

- `migrations/000001_create_users.up.sql`
- `migrations/000002_create_driver_profiles_and_loads.up.sql`
- `migrations/000003_add_dispatch_columns.up.sql`
- `migrations/000004_create_tracking_events.up.sql`
- `migrations/000005_create_payments.up.sql`

### Run The API

```bash
go run ./cmd/server
```

Health check:

```bash
curl http://localhost:8080/health
```

For a request-by-request localhost flow, see [http.go](/c:/Users/HP%20ENVY/Desktop/Projects/safiri-logistics/http.go).

## Scalability Notes

This system is intentionally not over-distributed yet, but the likely bottlenecks are already clear.

Near-term bottlenecks:

- candidate search for matching as driver count grows
- assignment races if multiple app instances match the same load
- tracking fan-out under many concurrent viewers
- payment callback reconciliation and retries

Recommended next upgrades:

- Redis geo indexes for nearby-driver candidate lookup
- Redis locks or DB-safe assignment guards
- Kafka or outbox pattern for assignment, tracking, and payment events
- Swagger or OpenAPI for external integration
- ETA prediction and AI-assisted dispatch ranking

## AI Opportunities

High-value AI additions for this domain:

- dispatch copilot with ranked driver recommendations and reason codes
- ETA prediction using lane history, borders, weather, and traffic
- anomaly detection for route deviation, idle time, and delivery risk
- pricing suggestions from lane demand and seasonality
- operations assistant queries such as:
  - `Which high-priority loads risk missing pickup tomorrow?`
  - `Which drivers can move a 32-ton lowbed load from Mombasa?`

The key product lesson is simple: AI becomes useful only when operational data is clean. Good driver state, asset/equipment data, load constraints, and tracking quality matter more than model complexity early on.
