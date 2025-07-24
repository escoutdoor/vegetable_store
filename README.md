# Vegetable Store

A microservice system for vegetable store management, built on Go using gRPC for inter-service communication and a complete observability stack.

## Architecture Overview

The system consists of three independent microservices:

- **Vegetable Service** - Manages vegetable inventory and catalog
- **User Service** - Handles user authentication and profile management
- **Order Service** - Processes and manages customer orders

## Tech Stack

### Core Technologies
- **Language**: Go
- **API**: 
  - gRPC for inter-service communication
  - REST Gateway (grpc-gateway) for external HTTP APIs
- **Documentation**: Swagger 2.0
- **Database**: PostgreSQL (separate instance per service)
- **Containerization**: Docker & Docker Compose

### Observability Stack
- **Metrics**: Prometheus
- **Tracing**: Jaeger with OpenTelemetry Collector
- **Logging**: Grafana Loki & Alloy
- **Visualization**: Grafana Dashboards
