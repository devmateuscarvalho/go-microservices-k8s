# Go Microservices with Kubernetes

## Project Overview
This project demonstrates a microservices architecture using Go, Kubernetes, and Docker.

## Services
- **Auth Service**: Handles user authentication and authorization
- **API Service**: Main application service
- **PostgreSQL**: Database for storing application data

## Kubernetes Configuration
The project uses Kubernetes for container orchestration, with configurations for:
- Deployments
- Services
- Persistent Volume Claims
- Ingress

## Local Development

### Prerequisites
- Docker Desktop
- Kubernetes enabled
- Go 1.20+

### Running the Project
1. Build Docker images for each service
2. Apply Kubernetes configurations:
   ```bash
   kubectl apply -f k8s/
   ```
3. Access services via Ingress at `localhost`

## Endpoints
- Auth Service: `/auth/*`
- API Service: `/api/*`

## Technologies
- Go
- Kubernetes
- Docker
- PostgreSQL
- NGINX Ingress Controller
