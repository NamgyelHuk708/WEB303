# 🚀 Practical 3: Full-Stack Microservices with gRPC, Databases, and Service Discovery

[![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![PostgreSQL](https://img.shields.io/badge/postgresql-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Consul](https://img.shields.io/badge/consul-%23F24C53.svg?style=for-the-badge&logo=consul&logoColor=white)](https://www.consul.io/)

This project demonstrates a complete microservices ecosystem built with modern cloud-native technologies including Go, gRPC, PostgreSQL, and Consul service discovery with dynamic routing capabilities.

---

## 📎 Repository

[GitHub Repository](https://github.com/Kinley-pal8/Web303_p3)

---

## 🏗️ Architecture Overview

The system implements a distributed microservices architecture with the following components:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   API Gateway   │◄──►│ Consul Discovery │◄──►│  Microservices  │
│   (Port 8080)   │    │   (Port 8500)    │    │                 │
│                 │    │                  │    │ Users Service   │
│ 1. Query Consul │    │ Service Registry │    │ (Port 50051)    │
│ 2. Get Addresses│    │                  │    │                 │
│ 3. Connect to   │    │ Health Checks    │    │ Products Service│
│    Services     │    │                  │    │ (Port 50052)    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                               │                 │
┌─────────────────┐    ┌──────────────────┐    │                 │
│   PostgreSQL    │◄───│   PostgreSQL     │◄───┘                 │
│   Users DB      │    │   Products DB    │                      │
│   (Port 5432)   │    │   (Port 5433)    │                      │
└─────────────────┘    └──────────────────┘                      │
```

### 🔧 Core Components

| Component             | Technology       | Port  | Purpose                                      |
| --------------------- | ---------------- | ----- | -------------------------------------------- |
| **API Gateway**       | Go + Gorilla Mux | 8080  | HTTP REST API with dynamic service discovery |
| **Users Service**     | Go + gRPC + GORM | 50051 | User management microservice                 |
| **Products Service**  | Go + gRPC + GORM | 50052 | Product management microservice              |
| **Consul**            | HashiCorp Consul | 8500  | Service discovery and health monitoring      |
| **Users Database**    | PostgreSQL 13    | 5432  | Isolated user data persistence               |
| **Products Database** | PostgreSQL 13    | 5433  | Isolated product data persistence            |

## ✨ Key Features

- 🔄 **Dynamic Service Discovery**: API Gateway queries Consul on each request for real-time service location
- 🚀 **gRPC Communication**: High-performance binary protocol for inter-service communication
- 🛡️ **Database Isolation**: Each microservice maintains its own PostgreSQL database
- 🔗 **Protocol Buffers**: Strongly-typed service contracts with automatic code generation
- ⚡ **Concurrent Processing**: Parallel gRPC calls for efficient data aggregation
- 🐳 **Full Containerization**: Docker-based deployment with multi-stage builds
- 📊 **Health Monitoring**: Automatic service registration and health checking

## 🚀 Quick Start

### Prerequisites

Ensure you have the following installed:

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
- [Go 1.18+](https://golang.org/dl/) (for local development)
- [Protocol Buffers Compiler](https://grpc.io/docs/protoc-installation/)
- [Buf CLI](https://docs.buf.build/installation) (for proto generation)

### 🛠️ Installation & Setup

1. **Clone and navigate to the project**:

   ```bash
   git clone https://github.com/Kinley-pal8/Web303_p3.git
   cd practical-three
   ```

2. **Generate Protocol Buffers** (if needed):

   ```bash
   buf generate
   ```

3. **Build and run the entire stack**:

   ```bash
   docker-compose up --build
   ```

4. **Verify services are running**:
   ```bash
   docker-compose ps
   ```

### 🔍 Service Health Check

Visit the Consul UI at [http://localhost:8500](http://localhost:8500) to monitor service health and discovery.

## 📚 API Documentation

### 👥 User Management

#### Create User

```bash
curl -X POST -H "Content-Type: application/json" \
     -d '{"name": "Jane Doe", "email": "jane.doe@example.com"}' \
     http://localhost:8080/api/users
```

**Response:**

```json
{ "id": "1", "name": "Jane Doe", "email": "jane.doe@example.com" }
```

#### Get User by ID

```bash
curl http://localhost:8080/api/users/1
```

### 🛍️ Product Management

#### Create Product

```bash
curl -X POST -H "Content-Type: application/json" \
     -d '{"name": "Laptop", "price": 1200.50}' \
     http://localhost:8080/api/products
```

**Response:**

```json
{ "id": "1", "name": "Laptop", "price": 1200.5 }
```

#### Get Product by ID

```bash
curl http://localhost:8080/api/products/1
```

### 🔗 Data Aggregation

#### Get Combined Purchase Data

```bash
curl http://localhost:8080/api/purchases/user/1/product/1
```

**Response:**

```json
{
  "user": {
    "id": "1",
    "name": "Jane Doe",
    "email": "jane.doe@example.com"
  },
  "product": {
    "id": "1",
    "name": "Laptop",
    "price": 1200.5
  }
}
```

## 🧪 Live API Testing Results

The following screenshot demonstrates the successful implementation and testing of all API endpoints:

![alt text](<Screenshot 2025-08-27 at 2.52.10 PM.png>)

As shown above, the system successfully:

- ✅ Retrieves user data via dynamic service discovery
- ✅ Fetches product information through Consul routing
- ✅ Aggregates data from multiple microservices in real-time

## 📁 Project Structure

```
practical-three/
├── 🌐 api-gateway/                 # HTTP REST API Gateway
│   ├── main.go                     # Dynamic service discovery logic
│   ├── Dockerfile                  # Multi-stage build configuration
│   ├── go.mod                      # Go module dependencies
│   └── go.sum                      # Dependency checksums
├── 🔧 services/
│   ├── 👥 users-service/           # User management microservice
│   │   ├── main.go                 # gRPC server with database integration
│   │   ├── Dockerfile              # Service containerization
│   │   ├── go.mod                  # Service-specific dependencies
│   │   └── go.sum
│   └── 🛍️ products-service/        # Product management microservice
│       ├── main.go                 # gRPC server with database integration
│       ├── Dockerfile              # Service containerization
│       ├── go.mod                  # Service-specific dependencies
│       └── go.sum
├── 📡 proto/                       # Protocol Buffer definitions
│   ├── users.proto                 # User service contract
│   ├── products.proto              # Product service contract
│   └── gen/proto/                  # Generated Go code
│       ├── go.mod                  # Proto module
│       ├── users.pb.go             # User protobuf bindings
│       ├── users_grpc.pb.go        # User gRPC client/server
│       ├── products.pb.go          # Product protobuf bindings
│       └── products_grpc.pb.go     # Product gRPC client/server
├── 🐳 docker-compose.yml           # Multi-service orchestration
├── 📋 buf.yaml                     # Buf configuration
├── ⚙️ buf.gen.yaml                 # Code generation settings
└── 📖 README.md                    # This documentation
```

## 🔬 Advanced Features

### Dynamic Service Discovery Flow

1. **Request Reception**: API Gateway receives HTTP request
2. **Consul Query**: Gateway queries Consul for service location
3. **Address Resolution**: Consul returns current service address (container:port)
4. **gRPC Connection**: Gateway establishes connection to discovered service
5. **Business Logic**: Service processes request with database interaction
6. **Response Aggregation**: Gateway combines multiple service responses (if needed)
7. **Client Response**: Final result returned to client

### Database Architecture

Each microservice maintains complete data isolation:

- **Users Service**: Independent PostgreSQL instance on port 5432
- **Products Service**: Separate PostgreSQL instance on port 5433
- **GORM Integration**: Automatic schema migration and ORM mapping
- **Connection Resilience**: Retry logic for database connectivity

## 🚀 Development Workflow

### Local Development Setup

1. **Start infrastructure services**:

   ```bash
   docker-compose up consul users-db products-db
   ```

2. **Run services locally** (in separate terminals):

   ```bash
   # Users Service
   cd services/users-service && go run main.go

   # Products Service
   cd services/products-service && go run main.go

   # API Gateway
   cd api-gateway && go run main.go
   ```

### Protocol Buffer Regeneration

```bash
# Install dependencies
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Regenerate code
buf generate
```

### Service Testing

```bash
# Test individual services
curl http://localhost:8080/api/users/1
curl http://localhost:8080/api/products/1

# Test service aggregation
curl http://localhost:8080/api/purchases/user/1/product/1

# Monitor service discovery
docker-compose logs api-gateway | grep "Discovering"
```

## 🎯 Learning Objectives Achieved

| Objective | Implementation          | Evidence                                        |
| --------- | ----------------------- | ----------------------------------------------- |
| **LO2**   | gRPC & Protocol Buffers | Efficient binary communication between services |
| **LO4**   | Data Persistence        | GORM integration with PostgreSQL                |
| **LO8**   | Service Discovery       | Dynamic Consul-based routing                    |
| **LO10**  | Distributed Systems     | Microservices architecture with isolation       |

## 📊 Monitoring & Observability

### Available Monitoring Tools

- **Consul UI**: [http://localhost:8500](http://localhost:8500) - Service health dashboard
- **Service Logs**: `docker-compose logs <service-name>` - Individual service monitoring
- **API Gateway Logs**: Real-time service discovery tracking
- **Database Monitoring**: PostgreSQL connection and query logging

### Health Check Endpoints

All services automatically register with Consul and provide health status information accessible through the Consul API and web interface.

## 🔧 Configuration

### Environment Variables

| Variable         | Default                | Description           |
| ---------------- | ---------------------- | --------------------- |
| `CONSUL_ADDRESS` | `consul:8500`          | Consul server address |
| `DB_HOST`        | `users-db/products-db` | Database host         |
| `DB_USER`        | `user`                 | Database username     |
| `DB_PASSWORD`    | `password`             | Database password     |
| `GRPC_PORT`      | `50051/50052`          | gRPC service ports    |

### Docker Compose Configuration

The `docker-compose.yml` orchestrates all services with proper networking, dependencies, and health checks.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

**Built with ❤️ using Go, gRPC, PostgreSQL, and Consul**
