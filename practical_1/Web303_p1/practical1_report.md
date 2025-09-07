# Practical 1 Report: Microservices Development with Go, gRPC, and Docker

## Executive Summary

This report details the completion of Practical 1 for the WEB303 module, focusing on setting up a foundational development environment for microservices and implementing inter-service communication using Go, gRPC, Protocol Buffers, and Docker Compose. The practical involved creating two microservices—a `time-service` that provides current time data and a `greeter-service` that generates personalized greetings while calling the time-service—and orchestrating them in a containerized environment. The implementation successfully demonstrated efficient inter-service communication via gRPC, with verification through Docker Compose logs and gRPC client testing.

The project aligns with module learning outcomes, including designing microservices with gRPC (Outcome 2), understanding container orchestration (Outcome 6), and explaining microservices concepts (Outcome 1).

## Objectives

The primary objectives of Practical 1 were:

- **Environment Setup:** Install and configure Go, Protocol Buffers, gRPC tools, and Docker to support microservice development.
- **Service Design:** Define service contracts using Protocol Buffers for a time service and a greeter service.
- **Implementation:** Build two Go-based microservices that communicate via gRPC.
- **Containerization:** Package services into Docker containers and orchestrate them using Docker Compose.
- **Verification:** Test inter-service communication and ensure services run correctly in a multi-container setup.

## Methodology

### 1. Environment Setup
- **Go Installation:** Downloaded and installed Go 1.24.2. Verified installation with:
  ```bash
  go version
  go env
  ```
- **Protocol Buffers and gRPC Tools:** Installed `protoc` compiler and Go plugins. Ran:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
  ```
  Updated PATH with:
  ```bash
  export PATH="$PATH:$(go env GOPATH)/bin"
  ```
- **Docker:** Installed Docker Desktop. Verified with:
  ```bash
  docker run hello-world
  ```

### 2. Project Structure and Service Contracts
- Created the directory structure:
  ```bash
  mkdir practical-one
  cd practical-one
  mkdir -p proto/gen
  mkdir greeter-service
  mkdir time-service
  ```
- Defined `proto/time.proto` with `TimeService` for retrieving current time.
- Defined `proto/greeter.proto` with `GreeterService` for generating greetings.
- Generated Go code:
  ```bash
  protoc --go_out=./proto/gen --go_opt=paths=source_relative \
      --go-grpc_out=./proto/gen --go-grpc_opt=paths=source_relative \
      proto/*.proto
  ```

### 3. Microservice Implementation
- **Time-Service:** Navigated to `time-service`, initialized module, and got dependencies:
  ```bash
  cd time-service
  go mod init practical-one/time-service
  go get google.golang.org/grpc
  ```
  Implemented `main.go` as a gRPC server.
- **Greeter-Service:** Navigated to `greeter-service`, initialized module, and got dependencies:
  ```bash
  cd greeter-service
  go mod init practical-one/greeter-service
  go get google.golang.org/grpc
  ```
  Implemented `main.go` as a gRPC server and client.
- Used shared proto-generated code for consistent interfaces.

### 4. Containerization and Orchestration
- Created Dockerfiles for both services using multi-stage builds (Go builder and Alpine runtime).
- Configured `docker-compose.yml` to build and run services, with `greeter-service` depending on `time-service`.
- Ensured services communicate via Docker network using hostnames (`time-service:50052`).

### 5. Testing and Verification
- Built and ran services:
  ```bash
  docker-compose up --build
  ```
- Verified logs for successful startup and communication.
- Tested `greeter-service` endpoint:
  ```bash
  grpcurl -plaintext \
      -import-path ./proto -proto greeter.proto \
      -d '{"name": "WEB303 Student"}' \
      0.0.0.0:50051 greeter.GreeterService/SayHello
  ```

## Results

### Successful Build and Deployment
- Docker Compose built images for both services without errors.
- Services started correctly:
  - `time-service` listening on port 50052.
  - `greeter-service` listening on port 50051 and connecting to `time-service`.
- Logs confirmed inter-service communication within the Docker network.

### gRPC Testing Output
Using `grpcurl`:
```
grpcurl -plaintext \
    -import-path ./proto -proto greeter.proto \
    -d '{"name": "WEB303 Student"}' \
    0.0.0.0:50051 greeter.GreeterService/SayHello
```
**Response:**
```json
{
  "message": "Hello WEB303 Student! The current time is 2025-09-07T15:45:15Z"
}
```
This demonstrates successful gRPC call flow: `greeter-service` received the request, called `time-service` for the current time, and returned a combined response.

### Docker Compose Logs Excerpt
```
Creating web303_p1-main_time-service_1 ... done
Creating web303_p1-main_greeter-service_1 ... done
Attaching to web303_p1-main_time-service_1, web303_p1-main_greeter-service_1
time-service_1     | 2025/09/07 15:45:15 Time service listening at [::]:50052
greeter-service_1  | 2025/09/07 15:45:15 Greeter service listening at [::]:50051
```

### Explanation of Inter-Service Communication
The `greeter-service` is able to find and call the `time-service` within the Docker network through the following mechanism:
- Docker Compose creates a default network for the services, allowing them to communicate using their service names as hostnames.
- In the `greeter-service` code, a gRPC client dials `time-service:50052`, where `time-service` is the hostname defined in `docker-compose.yml`.
- The `depends_on` directive ensures `time-service` starts before `greeter-service`, guaranteeing availability.
- This internal networking enables secure, efficient communication without exposing ports externally, demonstrating microservices orchestration in a containerized environment.

## Discussion

### Challenges and Solutions
- **Module Path Conflicts:** Initial import errors due to mismatched Go module names. Resolved by consolidating to a single root module (`practical-one`) and updating Dockerfiles to copy root `go.mod`.
- **Service Code Errors:** `greeter-service` initially had incorrect implementation (duplicated time-service code). Fixed by implementing proper `SayHello` method and client connection to `time-service`.
- **Proto Generation:** Ensured `protoc` commands used correct paths and options for source-relative generation.
- **Docker Networking:** Verified service discovery via hostnames in `docker-compose.yml` for seamless inter-service calls.

### Key Learnings
- gRPC enables efficient, type-safe communication between microservices.
- Docker Compose simplifies multi-container orchestration, mimicking production environments.
- Proper module management is crucial for shared code in Go projects.
- Containerization ensures portability and isolation.

### Alignment with Learning Outcomes
- **Outcome 2:** Demonstrated gRPC and Protocol Buffers for service contracts and communication.
- **Outcome 6:** Used Docker Compose as a stepping stone to Kubernetes-style orchestration.
- **Outcome 1:** Illustrated microservices benefits (modularity, scalability) and trade-offs (network complexity).

## Conclusion

Practical 1 was successfully completed, resulting in a functional microservices application with gRPC-based communication, containerized and orchestrated via Docker Compose. The setup provides a solid foundation for further exploration of microservices architectures, serverless computing, and advanced deployment strategies.

The implementation meets all specified requirements, with verified inter-service communication and external API testing. Future enhancements could include adding authentication, error handling, or scaling with Kubernetes.

## Appendices

- **Code Repository:** [GitHub Link] (Replace with actual repository URL upon submission).
- **Full Docker Compose Logs:** Available in the repository under `logs/`.
- **Proto Files:** Located in `proto/` directory.
- **Service Code:** In `time-service/` and `greeter-service/` directories.

