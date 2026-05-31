# Run Factory

**Run Factory** is an open-source platform for orchestrating, visualizing, and scaling millions of robots working together in a factory.

It provides the infrastructure layer for building large-scale, distributed robotic systems — enabling real-time coordination, intelligent task allocation, and full-system visibility through digital twin simulation.

---

## 🚀 Overview

Modern factories are no longer collections of isolated machines — they are **distributed robotic systems**.

Traditional centralized control architectures do not scale.

Run Factory is designed from the ground up to support:

* Massive scale (from 1 → millions of robots)
* Real-time coordination
* Distributed intelligence
* Fault-tolerant operation

---

## 🧠 Core Principles

* **Edge Autonomy**
  Robots operate independently and remain functional even under network disruption.

* **Event-Driven Architecture**
  All interactions are based on real-time events, not polling or blocking calls.

* **Horizontal Scalability**
  Every component is designed to scale horizontally.

* **Zone-based Coordination**
  Robots are grouped into zones for efficient orchestration at scale.

* **Digital Twin First**
  The system is designed to be visualized, simulated, and understood in real time.

---

## ⚙️ Core Features

### 🕸️ Swarm Orchestration

Coordinate thousands to millions of robots using distributed scheduling and task allocation.

### 📡 Real-time Communication

Low-latency messaging via pub/sub protocols (MQTT, DDS).

### 🧭 Task Scheduling Engine

Assign, balance, and optimize tasks across robots dynamically.

### 🧊 Digital Twin Visualization

Monitor your entire factory in real time with a visual interface.

### 🧱 Zone-based Scaling

Divide large factories into logical zones for better performance and control.

### 🛡️ Fault Tolerance

Robots continue operating with degraded connectivity or partial system failure.

---

## 🏗️ Architecture

Run Factory is built on a layered architecture:

### 1. Edge Layer (Robots)

* ROS2 / micro-runtime
* Local decision making
* Sensor integration

### 2. Communication Layer

* MQTT / DDS (pub/sub)
* Event streaming backbone

### 3. Orchestration Layer

* Distributed scheduler
* Task allocation engine
* Coordination logic

### 4. Visualization Layer

* Real-time dashboard
* Digital twin simulation
* System monitoring

---

## 📊 Example Flow

1. A task is created (e.g., move item from A → B)
2. Orchestrator assigns the task to a robot
3. Robot executes locally with autonomy
4. Robot streams state updates in real time
5. System visualizes and optimizes continuously

---

## 🌍 Use Cases

* Autonomous warehouses
* Smart factories
* Robot farms
* Multi-robot research systems
* Industrial automation platforms

---

## 🔧 Tech Stack (Recommended)

* **Robot Layer**: ROS2, micro-ROS
* **Messaging**: MQTT, DDS
* **Streaming**: Kafka / NATS
* **Backend**: Go / Python
* **Database**: PostgreSQL, Time-series DB
* **Frontend**: React + WebGL (Three.js)

---

## 🧪 Getting Started

```bash
# Clone repository
git clone https://github.com/your-org/run-factory.git

cd run-factory

# Start backend, dashboard, and local support services
make up
```

For common local workflows:

```bash
make test  # backend test/build + dashboard production build
make lint  # dashboard lint
make down  # stop and remove local stack
```

This brings up:

* backend API at `http://localhost:8080/health`
* dashboard at `http://localhost:4173`
* PostgreSQL at `localhost:5432`
* Mosquitto MQTT/WebSocket broker at `localhost:1883` and `localhost:9001`
* NATS with JetStream at `localhost:4222`

> The bundled Mosquitto config enables anonymous local access for development only and should not be reused in production.

---

## 🛣️ Roadmap

* [ ] MVP orchestration engine
* [ ] Basic robot communication (MQTT)
* [ ] Real-time dashboard
* [ ] Multi-zone coordination
* [ ] Digital twin simulation
* [ ] AI-based task optimization
* [ ] Large-scale deployment support (100k+ robots)

---

## 🤝 Contributing

We welcome contributions from developers, roboticists, and system engineers.

Please open issues or submit pull requests to help improve the platform.

---

## 📜 License

MIT License

---

## 🌌 Vision

Run Factory aims to become the **infrastructure layer for physical AI systems** —
similar to what Kubernetes is for cloud-native computing.

From a single robot to millions,
Run Factory enables machines to operate as a unified, intelligent system.

---

> Build the future of autonomous factories.
