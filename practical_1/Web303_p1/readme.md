# Practical 1: Microservices with Go and Docker

## Overview

This practical demonstrates the implementation of two microservices in Go: a Greeter Service and a Time Service. The services communicate using gRPC and are containerized using Docker. Protocol Buffers are used for service definitions and message serialization.

## Structure

```
practical_1/
├── docker-compose.yml
├── greeter-service/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── proto/
│   ├── greeter.proto
│   ├── time.proto
│   └── gen/
│       └── proto/
│           ├── greeter_grpc.pb.go
│           ├── greeter.pb.go
│           ├── time_grpc.pb.go
│           └── time.pb.go
├── time-service/
    ├── Dockerfile
    ├── go.mod
    ├── go.sum
    └── main.go
```

## Services

### Greeter Service

- Responds to greeting requests via gRPC.
- Source: `greeter-service/main.go`

### Time Service

- Provides the current time via gRPC.
- Source: `time-service/main.go`

## Protocol Buffers

- Service definitions: `proto/greeter.proto`, `proto/time.proto`
- Generated Go code: `proto/gen/proto/`

## Running the Services

1. **Build and Start Containers**
   - Use Docker Compose:
     ```zsh
     cd practical_1
     docker-compose up --build
     ```
2. **Accessing the Services**
   - The services will be available on the ports defined in `docker-compose.yml`.
   - Use a gRPC client or tools like `grpcurl` to interact with the services.

## Notes

- Ensure Docker is installed and running.
- Protocol buffer files must be compiled before building the services. Use `protoc` with the Go plugin if you modify `.proto` files.

## Authors

- Kinley-pal8

## License

This project is for educational purposes.
