# Run Factory - 1 Month Execution Checklist (~50 tasks)

> Goal: break down delivery into small, executable issues/tasks.  
> Rule: completed items already in the repository are marked `[x]`.

## A) Baseline and project setup (10 tasks)
- [x] T01 - Keep root `docker-compose.yml` running backend, dashboard, PostgreSQL, Mosquitto, and NATS.
- [x] T02 - Keep root README "Getting Started" aligned with `docker compose up --build`.
- [x] T03 - Keep backend health endpoint `GET /health` returning `200 ok`.
- [x] T04 - Keep backend factory routes registered (`POST/GET /factories`, `GET /factories/{id}`, `PUT /factories/{id}/map`, robot add/remove).
- [x] T05 - Keep in-memory factory repository implementation for MVP baseline.
- [x] T06 - Keep backend modular boundaries (`domain/repository/service/transport` + `platform`) intact.
- [x] T07 - Keep backend validation commands green: `go test ./... && go build ./...`.
- [x] T08 - Keep dashboard validation commands green: `npm run lint && npm run build`.
- [x] T09 - Add root `.env.example` with explicit keys: `API_PORT`, `DATABASE_URL`, `MQTT_BROKER_URL`, `NATS_URL`, and dashboard API base URL.
- [ ] T10 - Add `Makefile` targets (`make up`, `make down`, `make test`, `make lint`) to standardize local workflow.

## B) Data model and persistence (12 tasks)
- [ ] T11 - Define PostgreSQL schema for `factories` table (id, name, map3d, created_at, updated_at).
- [ ] T12 - Define PostgreSQL schema for `robots` table (id, factory_id, name, position_x/y/z, status, updated_at).
- [ ] T13 - Define PostgreSQL schema for `zones` table (id, factory_id, name, polygon/bounds, metadata).
- [ ] T14 - Define PostgreSQL schema for `tasks` table (id, factory_id, robot_id, zone_id, type, payload, status, timestamps).
- [ ] T15 - Add DB migration runner and initial migration files.
- [ ] T16 - Implement PostgreSQL repository for factories (`create/list/get/update map`).
- [ ] T17 - Implement PostgreSQL repository for robot add/remove/list operations.
- [ ] T18 - Add repository transaction helper for multi-step updates.
- [ ] T19 - Add repository-level optimistic locking/version field for concurrent updates.
- [ ] T20 - Add indexes for hot queries (`factory_id`, `robot_id`, `status`, `updated_at`).
- [ ] T21 - Add data access tests for PostgreSQL repositories (success + conflict + not-found cases).
- [ ] T22 - Add `STORAGE_BACKEND` config (`memory|postgres`) with typed constants, plus fail-fast startup validation for unknown values or unavailable DB connection.

## C) Backend API expansion (10 tasks)
- [ ] T23 - Add endpoint `GET /factories/{factoryID}/robots` with pagination parameters.
- [ ] T24 - Add endpoint `PATCH /factories/{factoryID}/robots/{robotID}/position`.
- [ ] T25 - Add endpoint `PATCH /factories/{factoryID}/robots/{robotID}/status`.
- [ ] T26 - Add endpoint `POST /factories/{factoryID}/zones`.
- [ ] T27 - Add endpoint `GET /factories/{factoryID}/zones`.
- [ ] T28 - Add endpoint `PATCH /factories/{factoryID}/zones/{zoneID}`.
- [ ] T29 - Add endpoint `POST /factories/{factoryID}/tasks` (task creation).
- [ ] T30 - Add endpoint `POST /factories/{factoryID}/tasks/{taskID}/assign` (assign robot).
- [ ] T31 - Add endpoint `PATCH /factories/{factoryID}/tasks/{taskID}/status` (queued/running/done/failed).
- [ ] T32 - Add endpoint `GET /factories/{factoryID}/tasks` with filtering by status/robot/zone.

## D) Eventing and orchestration MVP (8 tasks)
- [ ] T33 - Define event schema for robot telemetry, task lifecycle, and zone events.
- [ ] T34 - Implement publisher for robot/task events to NATS subject namespace.
- [ ] T35 - Implement MQTT bridge topic structure for robot state updates.
- [ ] T36 - Add backend consumer to update robot state from incoming telemetry events.
- [ ] T37 - Add simple scheduler worker: assign queued tasks to available robots in same zone.
- [ ] T38 - Add retry and dead-letter handling for failed task executions.
- [ ] T39 - Add idempotency handling for task creation using `Idempotency-Key` header with key persistence + TTL (e.g., 1-2h by default, documented if increased) to reject duplicate creates.
- [ ] T40 - Add orchestration metrics (queue size, assignment latency, task success rate).

## E) Dashboard delivery plan (6 tasks)
- [x] T41 - Keep React/Vite dashboard scaffold running as frontend baseline.
- [ ] T42 - Replace template home page with factory operations landing page.
- [ ] T43 - Implement API client module for backend endpoints (`health/factories/robots/tasks`).
- [ ] T44 - Implement factory list + create form UI.
- [ ] T45 - Implement factory detail page (map, robot list, task list).
- [ ] T46 - Implement robot detail panel with live position/status updates.

## F) Quality, security, CI/CD, operations (4 tasks)
- [ ] T47 - Add GitHub Actions workflow for backend test/build on pull requests.
- [ ] T48 - Add GitHub Actions workflow for dashboard lint/build on pull requests.
- [ ] T49 - Add API authentication middleware (token-based) for write endpoints.
- [ ] T50 - Add production hardening checklist for Mosquitto/NATS/PostgreSQL credentials and network policies.
