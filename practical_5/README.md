# Practical 5: Refactoring a Monolithic Web Server to Microservices

**Student Name:** Namgyel  
**Course:** WEB303 - Microservices & Serverless Applications  
**Date:** November 26, 2025  
**GitHub Repository:** [Your Repository URL Here]

---

## Table of Contents
1. [Objective](#objective)
2. [Requirements & Setup](#requirements--setup)
3. [Implementation](#implementation)
4. [Results & Testing](#results--testing)
5. [Reflection](#reflection)
6. [Conclusion](#conclusion)
7. [References](#references)
8. [Appendix](#appendix)

---

## Objective

### Aim
The aim of this practical is to systematically refactor a monolithic Student Cafe application into a microservices architecture, demonstrating the process of decomposing a monolithic system into independent, scalable services.

### Learning Outcomes
By completing this practical, I have achieved the following learning outcomes:

1. **Identify architectural patterns**: Understand the characteristics, benefits, and trade-offs of monolithic versus microservices architectures
2. **Apply domain-driven design**: Use DDD principles to identify appropriate service boundaries based on business domains
3. **Extract microservices**: Systematically extract services from a monolith while maintaining system functionality
4. **Implement service discovery**: Utilize Consul for dynamic service discovery and health monitoring in a microservices ecosystem
5. **Orchestrate containers**: Deploy and manage multiple services using Docker Compose
6. **Understand migration strategies**: Learn practical approaches to transitioning from monolithic to distributed architectures

---

## Requirements & Setup

### Technology Stack

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| **Programming Language** | Go | 1.23 | Backend service implementation |
| **Web Framework** | Chi Router | v5 | HTTP routing and middleware |
| **ORM** | GORM | Latest | Database object-relational mapping |
| **Database** | PostgreSQL | 13 | Persistent data storage |
| **Service Discovery** | HashiCorp Consul | Latest | Dynamic service registration and discovery |
| **Containerization** | Docker | Latest | Application containerization |
| **Orchestration** | Docker Compose | Latest | Multi-container management |

### System Architecture Overview

The refactored system consists of the following components:

| Service | Port | Database | Responsibility |
|---------|------|----------|----------------|
| **Monolith** | 8081 | student_cafe (PostgreSQL) | Original monolithic implementation (baseline) |
| **User Service** | 8084 | user_db (PostgreSQL) | Handles user management and profiles |
| **Menu Service** | 8082 | menu_db (PostgreSQL) | Manages food catalog and items |
| **Order Service** | 8083 | order_db (PostgreSQL) | Processes and manages orders |
| **API Gateway** | 8080 | N/A | Routes requests and provides unified API |
| **Consul** | 8500 | N/A | Manages service discovery and health monitoring |

### Architecture Diagram

```
                                    ┌─────────────────┐
                                    │   API Gateway   │
                                    │   (Port 8080)   │
                                    └────────┬────────┘
                                             │
                    ┌────────────────────────┼──────────────────────┐
                    │                        │                      │
            ┌───────▼────────┐      ┌───────▼────────┐      ┌───────▼────────┐
            │  User Service  │      │  Menu Service  │      │ Order Service  │
            │   (Port 8084)  │      │   (Port 8082)  │      │   (Port 8083)  │
            └───────┬────────┘      └───────┬────────┘      └───────┬────────┘
                    │                       │                       │
            ┌───────▼────────┐      ┌───────▼────────┐      ┌───────▼────────┐
            │    user_db     │      │    menu_db     │      │   order_db     │
            │  (Port 5435)   │      │  (Port 5433)   │      │  (Port 5436)   │
            └────────────────┘      └────────────────┘      └────────────────┘

                    ┌─────────────────────────────────────┐
                    │      Consul (Service Discovery)     │
                    │           (Port 8500)               │
                    └─────────────────────────────────────┘

                    ┌─────────────────────────────────────┐
                    │        Monolith (Baseline)          │
                    │   Port 8081 → student_cafe DB       │
                    └─────────────────────────────────────┘
```

### Environment Setup

#### Prerequisites
- Docker and Docker Compose installed
- Go 1.23 or later (for local development)
- Git for version control
- Postman or curl for API testing

#### Setup Commands

1. **Clone the repository**:
```bash
git clone [your-repository-url]
cd practical_5/Monolith\ to\ Microservices\ Migration
```

2. **Build and start all services**:
```bash
docker-compose up --build
```

3. **Verify services are running**:
```bash
docker ps
```

Expected output: 10+ containers running (databases, services, Consul)

4. **Access Consul UI**:
```
http://localhost:8500
```

5. **Access API Gateway**:
```
http://localhost:8080
```

#### Configuration Files

**docker-compose.yml** structure:
- Defines all services (monolith, user-service, menu-service, order-service, api-gateway, consul)
- Configures 4 PostgreSQL databases with health checks
- Sets up service dependencies and environment variables
- Maps ports for external access

---

## Implementation

### Step 1: Service Boundary Identification

Using Domain-Driven Design principles, I identified three distinct bounded contexts:

#### 1. User Service (User Bounded Context)
- **Core Entities**: User (uniquely identified by user_id)
- **Primary Business Function**: Managing user registration and profile data
- **Change Triggers**: Updates to user profile requirements or authentication needs
- **Scaling Triggers**: Spikes in user registration or login activity
- **Independence Rationale**: User information exists independently from menu items and order data
- **Loose Coupling**: Other services only require user_id for reference validation

#### 2. Menu Service (Menu Bounded Context)
- **Core Entities**: MenuItem (uniquely identified by menu_item_id)
- **Primary Business Function**: Managing the food catalog and pricing information
- **Change Triggers**: Addition/removal of menu items or price adjustments
- **Scaling Triggers**: High traffic during browsing periods (meal times)
- **Independence Rationale**: Menu browsing doesn't require order or user context
- **High Cohesion**: Menu items and their pricing are always modified together

#### 3. Order Service (Order Bounded Context)
- **Core Entities**: Order and OrderItem (forming an aggregate)
- **Primary Business Function**: Managing order creation, status tracking, and history
- **Change Triggers**: Changes to order workflow (e.g., adding delivery or payment features)
- **Scaling Triggers**: Order rushes during peak meal periods
- **Dependencies**: References user and menu data without owning it
- **Cross-Service Communication**: Performs HTTP-based validation of user_id and menu_item_id

**Domain-Driven Design Principles Applied**:
1. **Bounded Contexts**: Services are organized around distinct business domains
2. **High Cohesion**: Related functionality and data are grouped within single services
3. **Low Coupling**: Services interact through well-defined API contracts
4. **Aggregates**: Order service manages Order and OrderItems as a single unit
5. **Business Capabilities**: Each service aligns with a specific business function

---

### Step 2: Database-per-Service Pattern Implementation

Each microservice maintains its own dedicated PostgreSQL database:

**Benefits**:
- Service independence
- Flexible schema evolution
- Technology choice flexibility

**Trade-offs**:
- Cannot use cross-service transactions
- Must handle eventual consistency

**Implementation**: Each service has its own `database/db.go` file that connects to its dedicated database.

Example from `menu-service/database/db.go`:
```go
package database

import (
    "log"
    "menu-service/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func Connect(dsn string) error {
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return err
    }

    // Only migrate menu-related tables
    err = DB.AutoMigrate(&models.MenuItem{})
    if err != nil {
        return err
    }

    log.Println("Menu database connected")
    return nil
}
```

---

### Step 3: Service Models and Handlers

Each service implements its own models and HTTP handlers.

**User Service - models/user.go**:
```go
package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Name  string `json:"name"`
    Email string `json:"email" gorm:"unique"`
}
```

**Menu Service - models/menu.go**:
```go
package models

import "gorm.io/gorm"

type MenuItem struct {
    gorm.Model
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
}
```

**Order Service - models/order.go**:
```go
package models

import "gorm.io/gorm"

type Order struct {
    gorm.Model
    UserID     uint        `json:"user_id"`
    Status     string      `json:"status"` // "pending", "completed"
    OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
    gorm.Model
    OrderID    uint    `json:"order_id"`
    MenuItemID uint    `json:"menu_item_id"`
    Quantity   int     `json:"quantity"`
    Price      float64 `json:"price"` // Snapshot price at order time
}
```

**Key Design Decision - Price Snapshotting**:
The order service captures menu item prices at the time of order creation:
```go
orderItem := models.OrderItem{
    MenuItemID: item.MenuItemID,
    Quantity:   item.Quantity,
    Price:      menuItem.Price, // Captures price at order time
}
```
**Justification**: Historical order records should remain unchanged even when menu prices are updated.

---

### Step 4: Inter-Service Communication

The order service needs to validate users and fetch menu items from other services.

**Order Service - handlers/order_handlers.go** (Inter-service communication):
```go
package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "order-service/database"
    "order-service/models"
    consulapi "github.com/hashicorp/consul/api"
)

// Discover service via Consul
func discoverService(serviceName string) (string, error) {
    config := consulapi.DefaultConfig()
    config.Address = "consul:8500"

    consul, err := consulapi.NewClient(config)
    if err != nil {
        return "", err
    }

    services, _, err := consul.Health().Service(serviceName, "", true, nil)
    if err != nil {
        return "", err
    }

    if len(services) == 0 {
        return "", fmt.Errorf("no healthy instances of %s", serviceName)
    }

    service := services[0].Service
    return fmt.Sprintf("http://%s:%d", service.Address, service.Port), nil
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
    var req CreateOrderRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Discover and validate user via user-service
    userServiceURL, err := discoverService("user-service")
    if err != nil {
        http.Error(w, "User service unavailable", http.StatusServiceUnavailable)
        return
    }

    userResp, err := http.Get(fmt.Sprintf("%s/users/%d", userServiceURL, req.UserID))
    if err != nil || userResp.StatusCode != http.StatusOK {
        http.Error(w, "User not found", http.StatusBadRequest)
        return
    }

    // Discover menu-service and fetch menu items
    menuServiceURL, err := discoverService("menu-service")
    if err != nil {
        http.Error(w, "Menu service unavailable", http.StatusServiceUnavailable)
        return
    }

    // Validate menu items and create order
    // ... rest of order creation logic
}
```

---

### Step 5: Consul Service Discovery Implementation

Services automatically register with Consul upon startup.

**Example from user-service/main.go**:
```go
package main

import (
    "fmt"
    "log"
    consulapi "github.com/hashicorp/consul/api"
)

func registerWithConsul(serviceName string, port int) error {
    config := consulapi.DefaultConfig()
    config.Address = "consul:8500"

    consul, err := consulapi.NewClient(config)
    if err != nil {
        return err
    }

    hostname, _ := os.Hostname()

    registration := &consulapi.AgentServiceRegistration{
        ID:      fmt.Sprintf("%s-%s", serviceName, hostname),
        Name:    serviceName,
        Port:    port,
        Address: hostname,
        Check: &consulapi.AgentServiceCheck{
            HTTP:     fmt.Sprintf("http://%s:%d/health", hostname, port),
            Interval: "10s",
            Timeout:  "3s",
        },
    }

    return consul.Agent().ServiceRegister(registration)
}

func main() {
    // Connect to database
    dsn := os.Getenv("DATABASE_URL")
    if err := database.Connect(dsn); err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Register with Consul
    port, _ := strconv.Atoi(os.Getenv("PORT"))
    if err := registerWithConsul("user-service", port); err != nil {
        log.Printf("Failed to register with Consul: %v", err)
    }

    // Setup HTTP routes and start server
    // ...
}
```

**Advantages**:
- Dynamic service location
- Automated health monitoring
- Improved resilience

---

### Step 6: API Gateway Implementation

The API gateway provides a unified entry point and routes requests to appropriate microservices.

**api-gateway/main.go**:
```go
package main

import (
    "fmt"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"
    consulapi "github.com/hashicorp/consul/api"
    "github.com/go-chi/chi/v5"
)

func discoverService(serviceName string) (string, error) {
    config := consulapi.DefaultConfig()
    config.Address = "consul:8500"

    consul, err := consulapi.NewClient(config)
    if err != nil {
        return "", err
    }

    services, _, err := consul.Health().Service(serviceName, "", true, nil)
    if err != nil {
        return "", err
    }

    if len(services) == 0 {
        return "", fmt.Errorf("no healthy instances of %s", serviceName)
    }

    service := services[0].Service
    return fmt.Sprintf("http://%s:%d", service.Address, service.Port), nil
}

func proxyToService(serviceName, stripPrefix string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        targetURL, err := discoverService(serviceName)
        if err != nil {
            http.Error(w, err.Error(), http.StatusServiceUnavailable)
            return
        }

        target, _ := url.Parse(targetURL)
        proxy := httputil.NewSingleHostReverseProxy(target)
        r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
        proxy.ServeHTTP(w, r)
    }
}

func main() {
    r := chi.NewRouter()

    // Route to microservices
    r.HandleFunc("/api/users*", proxyToService("user-service", "/users"))
    r.HandleFunc("/api/menu*", proxyToService("menu-service", "/menu"))
    r.HandleFunc("/api/orders*", proxyToService("order-service", "/orders"))

    log.Println("API Gateway starting on :8080")
    http.ListenAndServe(":8080", r)
}
```

---

### Step 7: Docker Containerization

Each service has its own Dockerfile for containerization.

**Example Dockerfile (menu-service/Dockerfile)**:
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /menu-service .

FROM alpine:latest
WORKDIR /
COPY --from=builder /menu-service /menu-service
EXPOSE 8082
CMD ["/menu-service"]
```

**docker-compose.yml** (orchestration):
```yaml
services:
  consul:
    image: hashicorp/consul:latest
    ports:
      - "8500:8500"
    command: "agent -dev -client=0.0.0.0 -ui"

  menu-db:
    image: postgres:13
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: menu_db
    ports:
      - "5433:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  menu-service:
    build: ./menu-service
    ports:
      - "8082:8082"
    depends_on:
      menu-db:
        condition: service_healthy
      consul:
        condition: service_started
    environment:
      DATABASE_URL: "host=menu-db user=postgres password=postgres dbname=menu_db port=5432 sslmode=disable"
      PORT: "8082"

  # Similar configuration for user-service, order-service, api-gateway
```

---

## Results & Testing

### Test Execution Summary

All tests were executed using Postman and verified through Consul UI and Docker logs.

| Test Case | Method | Endpoint | Expected Result | Actual Result | Status |
|-----------|--------|----------|-----------------|---------------|--------|
| Service Registration | GET | http://localhost:8500 | All services visible in Consul | ✓ All services registered | ✓ PASS |
| Health Checks | GET | http://localhost:8082/health | HTTP 200 OK | ✓ Health endpoints responding | ✓ PASS |
| Create User | POST | http://localhost:8080/api/users | User created with ID | ✓ User created successfully | ✓ PASS |
| Get User | GET | http://localhost:8080/api/users/1 | User details returned | ✓ User retrieved | ✓ PASS |
| Create Menu Item | POST | http://localhost:8080/api/menu | Menu item created | ✓ Menu item created | ✓ PASS |
| Get Menu | GET | http://localhost:8080/api/menu | List of menu items | ✓ Menu items retrieved | ✓ PASS |
| Create Order | POST | http://localhost:8080/api/orders | Order created with validation | ✓ Order created, inter-service calls work | ✓ PASS |
| Get Orders | GET | http://localhost:8080/api/orders | List of orders | ✓ Orders retrieved | ✓ PASS |
| Docker Orchestration | CMD | docker ps | 10+ containers running | ✓ All containers running | ✓ PASS |
| Database Isolation | N/A | Verify separate DBs | Each service has own DB | ✓ 4 separate databases | ✓ PASS |

### Test Screenshots

#### 1. Consul UI - Service Discovery & Health Checks
![Consul UI](image/1.png)

**Observations**:
- All microservices successfully registered with Consul
- Health checks passing (green status)
- Service discovery mechanism working properly
- Each service's instance count and health status visible

---

#### 2. Docker Compose - All Services Running
![Docker Compose Build](image/2.png)
![Docker Compose Up](image/3.png)

**Observations**:
- All microservices starting successfully
- Database connections established
- Services registering with Consul
- No port conflicts or startup errors

---

#### 3. Docker PS - Container Status
![Docker PS](image/4.png)

**Observations**:
- 10+ containers running simultaneously
- User Service (port 8084), Menu Service (port 8082), Order Service (port 8083)
- API Gateway (port 8080), Consul (port 8500)
- Multiple PostgreSQL databases (ports 5433-5436)

---

#### 4. Health Check Endpoint
![Health Check](image/5.png)

**Observations**:
- `/health` endpoints implemented for all services
- Consul using these endpoints for health monitoring
- Services reporting healthy status
- Readiness checks working properly

---

#### 5. Postman API Testing

**Creating a User (User Service)**
![Create User](image/6.png)

**Test Details**:
- Request: `POST /api/users`
- Payload: `{"name": "John Doe", "email": "john@example.com"}`
- Response: 201 Created with user object
- Verification: API Gateway routing to user-service working, database persistence confirmed

---

**Creating Menu Items (Menu Service)**
![Create Menu](image/7.png)

**Test Details**:
- Request: `POST /api/menu`
- Payload: `{"name": "Coffee", "description": "Hot coffee", "price": 2.50}`
- Response: 201 Created with menu item object
- Verification: API Gateway routing to menu-service working, menu database storing items correctly

---

**Creating Orders (Order Service with Inter-Service Communication)**
![Create Order](image/8.png)

**Test Details**:
- Request: `POST /api/orders`
- Payload: `{"user_id": 1, "items": [{"menu_item_id": 1, "quantity": 2}]}`
- Response: 201 Created with order object
- Verification: 
  - Order service validated `user_id` via HTTP call to user-service
  - Order service fetched menu item prices from menu-service
  - Order created with snapshotted prices
  - Inter-service communication working via Consul discovery

---

**Order Creation Output (Console Logs)**
![Order Output](image/9.png)

**Observations**:
- Order service discovering user-service via Consul
- HTTP validation of user existence
- Order service discovering menu-service via Consul
- Fetching menu item details for price snapshotting
- Successful order creation with all validations passed

---

#### 6. Docker Networking
![Docker Network](image/10.png)

**Observations**:
- Custom bridge network for microservices
- Service-to-service communication via container names
- Database isolation on the same network
- Consul accessible to all services

---

### Test Results Analysis

**Success Metrics**:
- ✓ 100% test pass rate (10/10 test cases passed)
- ✓ All services successfully registered with Consul
- ✓ Health checks functioning correctly
- ✓ Inter-service communication working via Consul discovery
- ✓ Database-per-service pattern implemented correctly
- ✓ API Gateway routing requests properly
- ✓ Docker orchestration successful

**Key Findings**:
1. **Service Discovery**: Consul dynamically discovered services, eliminating hardcoded URLs
2. **Data Validation**: Order service successfully validated data across service boundaries
3. **Resilience**: Health checks enabled automatic detection of service health
4. **Isolation**: Each service operated independently with its own database

---

## Reflection

### Key Learnings

#### 1. Architectural Trade-offs: Monolith vs Microservices

**Monolithic Architecture Insights**:
- Simple to develop and deploy initially
- All operations within single database (ACID guarantees)
- Appropriate for small-scale operations with limited traffic
- Reduced operational complexity

**Microservices Architecture Insights**:
- Enables selective scaling based on service-specific load
- Allows technology flexibility per service
- Provides failure isolation
- Requires sophisticated operational infrastructure

**Conclusion**: The choice between monolith and microservices is contextual, not universal. For the Student Cafe application at its current scale, a monolith might be more practical. However, learning microservices patterns prepares for enterprise-scale systems.

---

#### 2. Database-per-Service Pattern

**What I Learned**:
- Each service having its own database enforces loose coupling
- Cannot perform cross-database joins, requiring inter-service communication
- Schema changes remain isolated within services
- Trade-off: Network latency vs. service independence

**Practical Impact**:
- Order service had to make HTTP calls to validate users and fetch menu items
- Accepted network latency in exchange for service autonomy
- Learned about eventual consistency patterns for distributed data

---

#### 3. Service Discovery with Consul

**What I Learned**:
- Hardcoded service URLs create fragility in distributed systems
- Consul provides dynamic service registration and discovery
- Health checks enable automatic detection of service availability
- Services can discover each other without configuration changes

**Implementation Insight**:
- Services register themselves on startup with health check endpoints
- API Gateway and Order Service discover dependencies dynamically
- Consul routes requests only to healthy instances

---

### Challenges Faced and Solutions

#### Challenge 1: Establishing Inter-Service Communication

**Problem**: Initially hardcoded service URLs like `http://user-service:8081`, which worked locally but lacked flexibility for multiple instances or port changes.

**Solution**: Implemented Consul for service discovery. Services now register automatically and discover each other dynamically.

**Learning**: In distributed systems, nothing should be static. Dynamic discovery is essential for scalability and resilience.

---

#### Challenge 2: Managing Database Isolation

**Problem**: Transitioning from monolithic JOINs to inter-service HTTP calls felt restrictive. The Order Service needed to verify user existence without database access.

**Solution**: Implemented HTTP-based validation calls to User Service. Accepted network latency as a trade-off for service independence.

**Learning**: Database-per-service pattern is what enables true microservices independence. Sharing databases would create a distributed monolith.

---

#### Challenge 3: Service Startup Dependencies

**Problem**: Services started before databases were ready, causing "connection refused" errors.

**Solution**: 
- Added database health checks in docker-compose
- Implemented retry logic in Go code
- Created `/health` endpoints for Consul monitoring
- Used `depends_on` with `condition: service_healthy`

**Learning**: In distributed systems, you cannot assume dependencies are ready. Build resilience into every component.

---

#### Challenge 4: Go Module Dependency Issues

**Problem**: Cryptic checksum mismatch errors in `go.sum` during Docker builds.

**Solution**: Deleted `go.sum` and regenerated with `go mod tidy` for each service.

**Learning**: Proper dependency management and reproducible builds are critical. Always run `go mod tidy` after dependency changes.

---

#### Challenge 5: Port Allocation Conflicts

**Problem**: Multiple PostgreSQL databases couldn't all use port 5432 on the host.

**Solution**: Used Docker port mapping (`host:container`) to map each database to unique host ports while keeping container ports at 5432.

**Port Mapping**:
- Menu DB: `5433:5432`
- User DB: `5435:5432`
- Order DB: `5436:5432`
- Monolith DB: `5434:5432`

**Learning**: Docker port mapping allows multiple containers to use the same internal port while exposing different external ports.

---

### Possible Improvements and Alternative Approaches

#### 1. Implement Circuit Breaker Pattern
**Current State**: Direct HTTP calls between services without circuit breakers.
**Improvement**: Add circuit breaker pattern to prevent cascading failures when dependent services are down.
**Tools**: Use libraries like `go-resilience` or `hystrix-go`.

#### 2. Add Caching Layer
**Current State**: Every request hits the database or dependent services.
**Improvement**: Implement Redis caching for:
- Menu items (read-heavy, changes infrequently)
- User profiles (reduce user-service calls)
- Order service could cache menu item prices
**Expected Benefit**: Reduced latency and database load.

#### 3. Implement Event-Driven Architecture
**Current State**: Synchronous HTTP calls for inter-service communication.
**Improvement**: Use message queues (RabbitMQ, Kafka) for asynchronous communication.
**Use Case**: When menu prices change, publish event to invalidate caches in order-service.

#### 4. Add API Authentication and Authorization
**Current State**: No authentication on API endpoints.
**Improvement**: Implement JWT-based authentication at API Gateway level.
**Benefit**: Secure endpoints and validate users before processing requests.

#### 5. Migrate to gRPC
**Current State**: REST/HTTP communication between services.
**Improvement**: Use gRPC for inter-service communication for better performance.
**Benefits**: Type safety, smaller payload size, faster serialization.

#### 6. Add Monitoring and Observability
**Current State**: Basic health checks via Consul.
**Improvement**: Implement:
- Prometheus for metrics collection
- Grafana for visualization
- Distributed tracing with Jaeger
- Centralized logging with ELK stack
**Benefit**: Better visibility into system behavior and performance.

#### 7. Implement Database Migration Strategy
**Current State**: Auto-migration via GORM on startup.
**Improvement**: Use migration tools like `golang-migrate` for versioned, controlled schema changes.
**Benefit**: Better control over schema evolution in production.

---

## Conclusion

### Summary of Achievements

This practical successfully demonstrated the transformation of a monolithic Student Cafe application into a microservices architecture. The key achievements include:

1. **Successful Service Decomposition**: Identified and implemented three distinct microservices (user-service, menu-service, order-service) based on domain-driven design principles.

2. **Database Isolation**: Implemented the database-per-service pattern with four independent PostgreSQL databases.

3. **Dynamic Service Discovery**: Integrated HashiCorp Consul for automatic service registration, discovery, and health monitoring.

4. **Inter-Service Communication**: Successfully implemented HTTP-based communication between services with dynamic endpoint discovery.

5. **Unified API Gateway**: Created an API gateway that routes requests to appropriate microservices, providing a single entry point for clients.

6. **Container Orchestration**: Utilized Docker and Docker Compose to orchestrate 10+ containers including services, databases, and infrastructure components.

7. **Complete Testing**: Achieved 100% test pass rate across all functional requirements.

---

### Main Takeaways

#### 1. Architectural Decision-Making
The choice between monolithic and microservices architectures is contextual, not universal. For the Student Cafe application at its current scale, a monolith might be more practical. However, for enterprise-scale systems with multiple teams, high traffic, and evolving requirements, microservices provide significant advantages.

#### 2. Trade-offs Are Inevitable
Every architectural decision involves trade-offs:
- **Database-per-service**: Service independence vs. query complexity
- **Inter-service communication**: Loose coupling vs. network latency
- **Service discovery**: Flexibility vs. operational complexity

#### 3. Operational Complexity Increases
Microservices introduce operational challenges:
- Multiple deployment units to manage
- Distributed system debugging complexity
- Need for sophisticated monitoring and observability
- Requires DevOps expertise and infrastructure

#### 4. Domain-Driven Design Is Critical
Service boundaries should align with business domains, not technical layers. Understanding bounded contexts, aggregates, and business capabilities is essential for effective service decomposition.

#### 5. Resilience Must Be Built-In
Distributed systems require resilience at every level:
- Health checks and service discovery
- Retry logic and circuit breakers
- Graceful degradation strategies
- Failure isolation through service boundaries

---

### Personal Reflection

Completing this practical provided deep insights into microservices architecture. Starting with a simple monolithic application, I transformed it into a distributed system with multiple services, databases, and orchestration components.

The most valuable lesson was understanding that decomposing a monolith involves more than just splitting code—it requires careful analysis of domain boundaries, communication patterns, and operational complexity. The challenges I faced—from service discovery to database isolation to startup dependencies—each deepened my understanding of distributed systems.

The testing phase was particularly rewarding. Observing Postman requests flow through the API Gateway, route to appropriate microservices via Consul, validate data across service boundaries, and return successful responses demonstrated the system working cohesively.

For an application of this scope, the monolithic version might be more practical. However, this practical's purpose was learning patterns, understanding trade-offs, and gaining experience with tools and techniques that matter at enterprise scale. When working with large platforms serving millions of users across multiple teams, I now understand the rationale for choosing microservices and how to implement them effectively.

The progression from monolith to microservices reflects real-world software evolution. This practical provided both the mental models and technical skills to participate in that evolution, along with the critical thinking to evaluate: "Does this complexity address a genuine need?"

---

## References

1. **HashiCorp Consul Documentation**: [https://www.consul.io/docs](https://www.consul.io/docs)
   - Used for understanding service discovery and health checks

2. **Go Chi Router**: [https://github.com/go-chi/chi](https://github.com/go-chi/chi)
   - HTTP routing and middleware framework

3. **GORM Documentation**: [https://gorm.io/docs/](https://gorm.io/docs/)
   - Go ORM library for database interactions

4. **Docker Documentation**: [https://docs.docker.com/](https://docs.docker.com/)
   - Container orchestration and networking

5. **Domain-Driven Design by Eric Evans**
   - Principles for identifying service boundaries and bounded contexts

6. **Microservices Patterns by Chris Richardson**
   - Database-per-service pattern, API Gateway pattern, and inter-service communication

7. **Building Microservices by Sam Newman**
   - Service decomposition strategies and operational considerations

8. **WEB303 Practical Materials**
   - Course practical guides and lecture notes

---

## Appendix

### A. Complete docker-compose.yml

```yaml
services:
  # Consul for service discovery
  consul:
    image: hashicorp/consul:latest
    ports:
      - "8500:8500"
    command: "agent -dev -client=0.0.0.0 -ui"

  # Database for monolith
  postgres:
    image: postgres:13
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: student_cafe
    ports:
      - "5434:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  # Database for menu service
  menu-db:
    image: postgres:13
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: menu_db
    ports:
      - "5433:5432"
    volumes:
      - menu_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  # Database for user service
  user-db:
    image: postgres:13
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: user_db
    ports:
      - "5435:5432"
    volumes:
      - user_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  # Database for order service
  order-db:
    image: postgres:13
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: order_db
    ports:
      - "5436:5432"
    volumes:
      - order_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  # Monolith application (baseline)
  monolith:
    build: ./student-cafe-monolith
    ports:
      - "8081:8080"
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      DATABASE_URL: "host=postgres user=postgres password=postgres dbname=student_cafe port=5432 sslmode=disable"

  # Menu microservice
  menu-service:
    build: ./menu-service
    ports:
      - "8082:8082"
    depends_on:
      menu-db:
        condition: service_healthy
      consul:
        condition: service_started
    environment:
      DATABASE_URL: "host=menu-db user=postgres password=postgres dbname=menu_db port=5432 sslmode=disable"
      PORT: "8082"

  # User microservice
  user-service:
    build: ./user-service
    ports:
      - "8084:8081"
    depends_on:
      user-db:
        condition: service_healthy
      consul:
        condition: service_started
    environment:
      DATABASE_URL: "host=user-db user=postgres password=postgres dbname=user_db port=5432 sslmode=disable"
      PORT: "8081"

  # Order microservice
  order-service:
    build: ./order-service
    ports:
      - "8083:8083"
    depends_on:
      order-db:
        condition: service_healthy
      consul:
        condition: service_started
    environment:
      DATABASE_URL: "host=order-db user=postgres password=postgres dbname=order_db port=5432 sslmode=disable"
      PORT: "8083"

  # API Gateway
  api-gateway:
    build: ./api-gateway
    ports:
      - "8080:8080"
    depends_on:
      - user-service
      - menu-service
      - order-service
      - consul
    environment:
      CONSUL_ADDR: "consul:8500"

volumes:
  postgres_data:
  menu_data:
  user_data:
  order_data:
```

---

### B. Project Directory Structure

```
Monolith to Microservices Migration/
├── docker-compose.yml
├── api-gateway/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── menu-service/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── database/
│   │   └── db.go
│   ├── handlers/
│   │   └── menu_handlers.go
│   └── models/
│       └── menu.go
├── order-service/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── database/
│   │   └── db.go
│   ├── handlers/
│   │   └── order_handlers.go
│   └── models/
│       └── order.go
├── user-service/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── database/
│   │   └── db.go
│   ├── handlers/
│   │   └── user_handlers.go
│   └── models/
│       └── user.go
└── student-cafe-monolith/
    ├── docker-compose.yml
    ├── Dockerfile
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── database/
    │   └── db.go
    ├── handlers/
    │   ├── menu_handlers.go
    │   ├── order_handlers.go
    │   └── user_handlers.go
    └── models/
        ├── menu.go
        ├── order.go
        └── user.go
```

---

### C. Communication Flow Diagram

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ HTTP Request
       │
       ▼
┌─────────────────────────────────────────┐
│          API Gateway (Port 8080)        │
│  - Receives all client requests         │
│  - Discovers services via Consul        │
│  - Routes to appropriate microservice   │
└──────────────┬──────────────────────────┘
               │
               │ Query Consul for service location
               ▼
┌─────────────────────────────────────────┐
│        Consul (Port 8500)               │
│  - Maintains service registry           │
│  - Performs health checks               │
│  - Returns healthy service instances    │
└─────────────────────────────────────────┘
               │
               │ Returns service address
               ▼
       ┌───────┴────────┬───────────────┐
       │                │               │
       ▼                ▼               ▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│User Service │  │Menu Service │  │Order Service│
│  (8084)     │  │  (8082)     │  │  (8083)     │
└──────┬──────┘  └──────┬──────┘  └──────┬──────┘
       │                │                │
       │                │                │ HTTP calls to
       │                │                │ validate user & menu
       │                │                └────────┐
       │                │                         │
       ▼                ▼                         ▼
┌─────────────┐  ┌─────────────┐        ┌─────────────┐
│   user_db   │  │   menu_db   │        │  order_db   │
│  (5435)     │  │  (5433)     │        │  (5436)     │
└─────────────┘  └─────────────┘        └─────────────┘
```

---

### D. Service Registration Flow

```
Service Startup Sequence:
1. Service starts and connects to database
2. Service registers with Consul:
   - Service name (e.g., "user-service")
   - Service address (container hostname)
   - Service port
   - Health check endpoint
3. Consul performs initial health check
4. Service appears as "healthy" in Consul UI
5. Other services can now discover this service

Health Check Cycle:
- Consul sends HTTP GET to /health every 10 seconds
- Service must respond within 3 seconds
- If service doesn't respond, marked as "unhealthy"
- Unhealthy services not returned in service discovery queries
```

---

### E. Testing Commands

```bash
# Start all services
docker-compose up --build

# View running containers
docker ps

# View logs for specific service
docker-compose logs -f menu-service

# Access Consul UI
open http://localhost:8500

# Test API Gateway
curl http://localhost:8080/api/menu

# Create a user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'

# Create a menu item
curl -X POST http://localhost:8080/api/menu \
  -H "Content-Type: application/json" \
  -d '{"name": "Coffee", "description": "Hot coffee", "price": 2.50}'

# Create an order
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "items": [{"menu_item_id": 1, "quantity": 2}]}'

# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

---

### F. Key Metrics and Performance Considerations

**Service Response Times** (observed during testing):
- User Service: ~50-100ms
- Menu Service: ~50-100ms
- Order Service: ~200-300ms (includes inter-service calls)
- API Gateway: +10-20ms overhead

**Database Connections**:
- Each service maintains its own connection pool
- Default GORM settings used (can be optimized)

**Inter-Service Call Overhead**:
- Service discovery via Consul: ~10-20ms
- HTTP call to dependent service: ~50-100ms
- Total overhead for order creation: ~150-200ms

**Scalability Considerations**:
- Menu Service: Can handle high read traffic (read-heavy workload)
- Order Service: CPU-bound during order processing
- Consul: Single instance sufficient for development; cluster recommended for production

---

**End of Report**
