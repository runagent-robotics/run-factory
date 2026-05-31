# Project TODO Checklist

## 1) Foundation & Setup
- [ ] Hoàn thiện hướng dẫn cài đặt local cho backend + dashboard + services qua Docker Compose
- [ ] Chuẩn hóa biến môi trường (`.env.example`) cho tất cả service
- [ ] Thiết lập workflow CI cơ bản cho backend và dashboard

## 2) Backend (runfactory-core)
- [ ] Thay repository in-memory bằng PostgreSQL persistence
- [ ] Thiết kế migration schema cho factory, robot, zone, task
- [ ] Bổ sung API quản lý task (create/assign/update status)
- [ ] Bổ sung API zone-based coordination
- [ ] Tích hợp event bus (MQTT/NATS) cho cập nhật trạng thái robot theo thời gian thực
- [ ] Bổ sung authn/authz cho API (token/API key/RBAC)
- [ ] Tăng cường observability (structured logging, metrics, tracing)

## 3) Dashboard (runfactory-dashboard)
- [ ] Thay thế giao diện Vite template bằng UI quản trị factory thực tế
- [ ] Kết nối dashboard với backend API (`/factories`, `/robots`, `/health`)
- [ ] Xây màn hình danh sách/chi tiết factory và robot
- [ ] Thêm realtime state view cho robot (WebSocket/MQTT bridge)
- [ ] Chuẩn hóa design system và state management

## 4) Digital Twin & Orchestration
- [ ] Xây bản MVP orchestration engine
- [ ] Hoàn thiện basic robot communication (MQTT)
- [ ] Triển khai multi-zone coordination
- [ ] Xây digital twin simulation realtime
- [ ] Nghiên cứu và tích hợp AI-based task optimization

## 5) Quality, Security, Operations
- [ ] Mở rộng unit/integration tests cho service, transport và repository
- [ ] Thiết lập E2E test cho luồng backend ↔ dashboard
- [ ] Hardening bảo mật cho broker/API (không dùng anonymous access ở production)
- [ ] Bổ sung rate limit + input validation nâng cao
- [ ] Viết runbook vận hành, backup/restore, và monitoring alerts
