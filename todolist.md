# Project TODO Checklist

## 1) Foundation & Setup
- [ ] Finalize local setup guide for backend + dashboard + support services via Docker Compose
- [ ] Standardize environment variables (`.env.example`) across all services
- [ ] Set up baseline CI workflows for backend and dashboard

## 2) Backend (runfactory-core)
- [ ] Replace in-memory repository with PostgreSQL persistence
- [ ] Design migration schema for factory, robot, zone, and task
- [ ] Add task management APIs (create/assign/update status)
- [ ] Add zone-based coordination APIs
- [ ] Integrate event bus (MQTT/NATS) for real-time robot state updates
- [ ] Add API authn/authz (token/API key/RBAC)
- [ ] Improve observability (structured logging, metrics, tracing)

## 3) Dashboard (runfactory-dashboard)
- [ ] Replace Vite starter UI with real factory administration screens
- [ ] Connect dashboard to backend APIs (`/factories`, `/robots`, `/health`)
- [ ] Build factory and robot list/detail pages
- [ ] Add real-time robot state view (WebSocket/MQTT bridge)
- [ ] Standardize design system and state management

## 4) Digital Twin & Orchestration
- [ ] Build MVP orchestration engine
- [ ] Complete basic robot communication (MQTT)
- [ ] Implement multi-zone coordination
- [ ] Build real-time digital twin simulation
- [ ] Research and integrate AI-based task optimization

## 5) Quality, Security, Operations
- [ ] Expand unit/integration tests for service, transport, and repository layers
- [ ] Set up E2E tests for backend ↔ dashboard flows
- [ ] Harden broker/API security (no anonymous access in production)
- [ ] Add rate limiting and stronger input validation
- [ ] Write operation runbooks, backup/restore plans, and monitoring alerts
