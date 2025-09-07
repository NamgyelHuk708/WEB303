# Web303_p2 Microservices Project

This project demonstrates a simple microservices architecture using Go. It consists of an API Gateway and two services: Products Service and Users Service.

## Project Structure

```
api-gateway/
    main.go
    go.mod
    go.sum
services/
    products-service/
        main.go
        go.mod
        go.sum
    users-service/
        main.go
        go.mod
        go.sum
```

## Components

### 1. API Gateway

- Entry point for client requests.
- Routes requests to the appropriate microservice.
- Handles request aggregation and response formatting.

### 2. Products Service

- Manages product-related operations (CRUD).
- Exposes RESTful endpoints for product data.

### 3. Users Service

- Manages user-related operations (CRUD).
- Exposes RESTful endpoints for user data.

## Getting Started

### Prerequisites

- Go 1.18 or higher

### Running the Services

1. **Start Products Service**
   ```zsh
   cd services/products-service
   go run main.go
   ```
2. **Start Users Service**
   ```zsh
   cd services/users-service
   go run main.go
   ```
3. **Start API Gateway**
   ```zsh
   cd api-gateway
   go run main.go
   ```

## API Endpoints

- Products Service: `http://localhost:<port>/products`
- Users Service: `http://localhost:<port>/users`
- API Gateway: `http://localhost:<port>/`

## Notes

- Each service runs independently and communicates via HTTP.
- Ports and further configuration can be set in each service's `main.go`.

## License

MIT
