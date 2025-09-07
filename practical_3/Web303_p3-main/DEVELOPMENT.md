# Development Guide - Practical 3

## Quick Start

1. **Start the entire system**:
   ```bash
   docker-compose up --build
   ```

2. **Test the API**:
   ```bash
   ./test-api.sh
   ```

3. **View Consul UI**: Open http://localhost:8500

## VS Code Extensions Installed

This project uses the following VS Code extensions for Protocol Buffer development:

- **Buf** (`bufbuild.vscode-buf`): Advanced Protocol Buffer support with linting and formatting
- **Protobuf** (`pbkit.vscode-pbkit`): Syntax highlighting and Go-to-definition for `.proto` files

## Development Workflow

### 1. Protocol Buffer Development
- Edit `.proto` files in the `proto/` directory
- Use Buf for code generation: `buf generate`
- The extensions provide syntax highlighting and validation

### 2. Service Development
- Each service has its own Go module
- Services auto-register with Consul on startup
- Use gRPC for inter-service communication

### 3. Testing
- API Gateway provides HTTP endpoints
- Use the test script for automated testing
- Check Consul UI for service health

## File Structure Explained

```
practical-three/
├── proto/                    # Protocol Buffer definitions
│   ├── users.proto          # User service contract
│   ├── products.proto       # Product service contract
│   ├── gen/                 # Generated Go code (auto-generated)
│   └── buf.yaml            # Buf configuration
├── services/                # Microservices
│   ├── users-service/       # User management service
│   └── products-service/    # Product management service
├── api-gateway/             # HTTP to gRPC gateway
├── docker-compose.yml       # Container orchestration
└── test-api.sh             # API testing script
```

## Common Commands

### Docker Management
```bash
# Start all services
docker-compose up --build

# Start in background
docker-compose up --build -d

# View logs
docker-compose logs <service-name>

# Stop all services
docker-compose down

# Remove volumes (reset databases)
docker-compose down -v
```

### Protocol Buffer Generation
```bash
# Generate Go code from proto files
buf generate

# Lint proto files
buf lint

# Format proto files
buf format -w
```

### Local Development
```bash
# Run individual services locally (after starting DBs)
cd services/users-service && go run main.go
cd services/products-service && go run main.go
cd api-gateway && go run main.go
```

## Troubleshooting

### Port Conflicts
- API Gateway: 8080
- Users Service: 50051
- Products Service: 50052
- Consul: 8500
- Users DB: 5432
- Products DB: 5433

### Service Discovery Issues
- Check Consul UI at http://localhost:8500
- Ensure services register successfully
- Verify network connectivity between containers

### Database Connection Issues
- Check database container logs: `docker-compose logs users-db`
- Verify connection strings in service main.go files
- Reset databases: `docker-compose down -v && docker-compose up --build`

### gRPC Generation Issues
- Ensure Buf is installed: `buf --version`
- Check buf.yaml and buf.gen.yaml configurations
- Re-run: `buf generate`

## Advanced Features

### Adding New Services
1. Create new `.proto` file
2. Update `buf.gen.yaml` if needed
3. Generate Go code: `buf generate`
4. Create service directory with main.go and Dockerfile
5. Add to docker-compose.yml
6. Update API Gateway if needed

### Custom Health Checks
Services can implement custom health check endpoints for Consul monitoring.

### Load Balancing
Consul supports load balancing across multiple instances of the same service.

### Security
- Add TLS/SSL certificates for production
- Implement authentication/authorization
- Use service mesh for advanced security features
