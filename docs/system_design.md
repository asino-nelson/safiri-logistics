# Safiri Logistics System Design

## 1. Problem Framing

Safiri is a heavy-goods logistics platform for trailers moving cargo such as machinery, timber, vehicles, and industrial equipment across African and global corridors. The near-term product goal is simple:

- customers post loads
- verified drivers come online with live location and equipment capacity
- the platform prices the job
- the closest eligible driver is matched automatically when possible
- admins can optimize dispatch for future loads in bulk
- customers and operations teams can track movements in real time

This document intentionally balances product ambition with pragmatic engineering. The current codebase should stay deployable as a single service while we introduce clear seams for Redis, Kafka, real-time tracking, and optimization.

## 2. Functional Requirements

- Users can register and sign in with JWT.
- Drivers must pass KYC before being eligible for assignments.
- Drivers can expose operational state:
  - online or offline
  - available or busy
  - current location
  - experience level
  - truck equipment type
  - max supported load weight
- Customers can post heavy-goods loads with:
  - pickup and dropoff locations
  - pickup schedule
  - weight
  - equipment requirement
  - cargo type
  - dispatch priority
- The system computes a freight quote and stores pricing metadata.
- The system attempts immediate matching to the closest eligible online driver.
- Admins can run batch dispatch optimization for future loads.
- Drivers can update load status and send location pings.
- Customers, drivers, and operators can subscribe to live tracking over WebSockets.
- Payments can be initiated via M-Pesa for now.

## 3. Non-Functional Requirements

- Low latency for single-load matching:
  - target p95 under 300 ms for in-region requests when candidate drivers are already warm in cache
- Horizontal scalability:
  - matching, tracking, and payments should not require sticky sessions
- Reliability:
  - dispatch decisions must be idempotent
  - payment initiation must be auditable
  - tracking events should be append-only
- Observability:
  - matching score, candidate count, and assignment reason should be explainable
- Fairness and business control:
  - nearest driver should not be the only signal
  - experience, capacity, and priority should influence assignments
- Extensibility:
  - keep seams for Redis caching, Kafka eventing, and later ML optimization

## 4. System Design Interview Style Review

### API Design

Existing:

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/me`
- `POST /api/v1/loads`
- `GET /api/v1/loads`
- `POST /api/v1/loads/:loadID/pick`
- `PATCH /api/v1/loads/:loadID/status`
- `POST /api/v1/drivers/kyc`
- `GET /api/v1/drivers/kyc/me`
- `PATCH /api/v1/drivers/:userID/kyc/review`

Added in this design direction:

- `PATCH /api/v1/drivers/me/operations`
  - updates online status, availability, live position, capacity, equipment type, and experience years
- `POST /api/v1/dispatch/optimize`
  - admin endpoint for batch matching loads scheduled for a day or time window
- `GET /api/v1/tracking/ws?load_id=:id`
  - WebSocket subscription for live updates
- `POST /api/v1/tracking/loads/:loadID/events`
  - driver location update endpoint
- `POST /api/v1/payments/loads/:loadID/mpesa-checkout`
  - initiate M-Pesa checkout

### High-Level Architecture

Phase-1 architecture stays intentionally simple:

- `API Service`
  - auth
  - order management
  - driver operations
  - matching engine
  - pricing engine
  - payment orchestration
  - tracking hub
- `PostgreSQL`
  - source of truth for users, driver profiles, loads, tracking events, payments
- `In-process event publisher`
  - keeps the local developer experience simple
- `WebSocket hub`
  - broadcasts tracking and assignment events to connected clients

Production-ready seams introduced now:

- `Redis`
  - hot cache for driver availability and geo-near candidate lookup
  - short-lived assignment locks
- `Kafka`
  - publishes `load.created`, `load.matched`, `tracking.updated`, `payment.initiated`
  - enables analytics, ETA prediction, billing, and alerting consumers later

### Deep Dive: Matching

The fastest valid assignment is usually:

1. filter eligible drivers
2. score candidates
3. assign atomically

Eligibility filters:

- driver KYC approved
- online and available
- supports required equipment type
- supports required load weight
- has current location

Candidate scoring:

- distance to pickup
- driver experience
- priority boost for urgent or high-complexity loads
- optional penalties for stale location, poor acceptance behavior, or overloaded region

Suggested heuristic:

`score = distance_km * 0.65 + experience_penalty * 0.20 + priority_penalty * 0.15`

Interpretation:

- lower score is better
- `experience_penalty` drops for more experienced drivers
- `priority_penalty` falls as load priority rises

This is intentionally not ML-first. It is explainable, debuggable, and good enough to ship.

### Deep Dive: Batch Dispatch Optimization

For tomorrow’s 50 loads:

1. fetch unassigned loads ordered by:
   - highest priority first
   - earliest pickup time first
   - heaviest or most complex cargo first
2. fetch eligible drivers once
3. greedily assign the best remaining driver per load
4. reserve a matched driver in-memory for the run so one driver is not double-booked
5. persist assignments and publish events

Why greedy first:

- simple
- understandable for dispatch teams
- cheap to run
- easy to replace later with OR-Tools or min-cost max-flow

### Deep Dive: Tracking

Tracking should separate write and read concerns:

- drivers send location events through HTTP
- API persists the event
- API broadcasts the event through WebSocket
- later Kafka consumers can build ETA, anomaly alerts, and route deviation detection

This gives reliable persistence plus low-latency fan-out without immediately requiring a streaming platform.

### Deep Dive: Pricing and M-Pesa

Pricing should be deterministic and explainable at first:

- base fee
- distance component
- weight component
- priority surcharge
- equipment surcharge for specialized trailers

Payments should be orchestrated, not embedded inside order creation:

- quote generated when load is created
- customer explicitly initiates checkout
- payment row records provider reference and state
- later callback reconciliation confirms paid status

### AI and Operations Strategy

The future AI advantage is not “AI everywhere.” It is AI where dispatch teams lose time:

- recommended driver ranking with reasons
- ETA prediction using historical corridor, border, weather, and traffic data
- anomaly detection for stalled trucks and route deviation
- automated KYC risk review assistance
- dynamic pricing suggestions from seasonality and lane demand
- natural-language ops copilot for:
  - “show me all high-priority loads at risk of SLA breach”
  - “which drivers can move 32-ton lowbed cargo from Mombasa tomorrow morning?”

Laminar Copilot appears, based on its help center, to focus on scheduling and managing carrier operations and to rely on structured contracts, assets, and driver data. One onboarding article explicitly says Copilot needs contracts, assets, and driver data to automatically populate schedules and assist with daily operations. That supports a strong product lesson for Safiri: dispatch quality depends on operational data quality, not just optimization logic.

## 5. Current Bottlenecks In The Codebase

- Single-process matching:
  - today all matching logic runs in-process, so scale-out could create duplicate assignment races unless guarded
- No assignment lock:
  - two concurrent matchers could assign the same driver or load
- Driver operational state is incomplete:
  - online status and fresh location are required for low-latency matching
- No event stream:
  - hard to add analytics, audit replay, and downstream consumers
- No payment workflow:
  - quoting without payment state creates a billing blind spot
- No tracking persistence:
  - live location without history would weaken dispute resolution and ETA analytics

## 6. Implementation Principles

- Keep the monolith, modularize the domain.
- Add interfaces where production infra will later sit.
- Prefer deterministic heuristics before machine learning.
- Make every assignment explainable.
- Persist every business-critical state transition.
- Use Redis and Kafka only where they remove a real bottleneck.

## 7. What To Improve Later

- Swagger or OpenAPI generation
- idempotency keys for payments and dispatch actions
- Redis geo indexes for candidate lookup
- Kafka-backed outbox pattern
- route optimization across multi-stop trips
- trailer subtype support:
  - flatbed
  - lowbed
  - car carrier
  - logging trailer
- admin dashboards and SLA alerting
- ETA prediction and AI-driven exception management
- driver quality and safety scoring

## 8. Practical Recommendation

The right next version is not a microservices rewrite. It is:

- one clean service
- stronger data model
- explainable heuristic matching
- persistent tracking
- payment orchestration
- explicit seams for Redis and Kafka

That gives a system that can serve early African freight corridors well, while staying credible for global expansion later.
