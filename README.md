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

## Local Development Setup

### Prerequisites
- Docker Desktop
- Go 1.20+
- kubectl
- Git

### Local Kubernetes Setup (Docker Desktop)

#### Windows
1. Install Docker Desktop
   - Download from: https://www.docker.com/products/docker-desktop
   - During installation, enable Kubernetes in settings
   
2. Verify Kubernetes Installation
   ```powershell
   kubectl version
   kubectl get nodes
   ```

#### macOS
1. Install Docker Desktop
   - Download from: https://www.docker.com/products/docker-desktop
   - Enable Kubernetes in preferences

2. Verify Kubernetes Installation
   ```bash
   kubectl version
   kubectl get nodes
   ```

#### Linux (Minikube)
1. Install Docker and Minikube
   ```bash
   # Install Docker
   sudo apt-get update
   sudo apt-get install docker.io

   # Install Minikube
   curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
   sudo install minikube-linux-amd64 /usr/local/bin/minikube

   # Start Minikube
   minikube start
   ```

### Running the Project Locally

1. Build Docker Images
   ```bash
   docker build -t auth-service:latest ./auth-service
   docker build -t api-service:latest ./api-service
   ```

2. Apply Kubernetes Configurations
   ```bash
   kubectl apply -f k8s/
   ```

3. Verify Deployment
   ```bash
   kubectl get pods
   kubectl get services
   kubectl get ingress
   ```

4. Access Services
   - Open `http://localhost/auth/health`
   - Open `http://localhost/api/health`

## AWS Deployment Guide

### Prerequisites
- AWS Account
- AWS CLI
- eksctl
- kubectl
- Docker

### Step-by-Step AWS EKS Deployment

1. Install AWS CLI and Configure
   ```bash
   # Install AWS CLI
   pip install awscli

   # Configure AWS Credentials
   aws configure
   ```

2. Install eksctl
   ```bash
   # macOS
   brew tap weaveworks/tap
   brew install weaveworks/tap/eksctl

   # Windows (using Chocolatey)
   choco install eksctl

   # Linux
   curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
   sudo mv /tmp/eksctl /usr/local/bin
   ```

3. Create EKS Cluster
   ```bash
   eksctl create cluster \
     --name go-microservices-cluster \
     --region us-east-1 \
     --nodegroup-name standard-workers \
     --node-type t3.medium \
     --nodes 3
   ```

4. Configure kubectl
   ```bash
   aws eks --region us-east-1 update-kubeconfig --name go-microservices-cluster
   ```

5. Install NGINX Ingress Controller
   ```bash
   kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.1.1/deploy/static/provider/aws/deploy.yaml
   ```

6. Build and Push Docker Images
   ```bash
   # Login to AWS ECR
   aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <your-account-id>.dkr.ecr.us-east-1.amazonaws.com

   # Build and Push Auth Service
   docker build -t <your-account-id>.dkr.ecr.us-east-1.amazonaws.com/auth-service:latest ./auth-service
   docker push <your-account-id>.dkr.ecr.us-east-1.amazonaws.com/auth-service:latest

   # Build and Push API Service
   docker build -t <your-account-id>.dkr.ecr.us-east-1.amazonaws.com/api-service:latest ./api-service
   docker push <your-account-id>.dkr.ecr.us-east-1.amazonaws.com/api-service:latest
   ```

7. Update Kubernetes Deployment Files
   - Replace image names in `k8s/auth-service.yaml` and `k8s/api-service.yaml` with ECR image URLs

8. Deploy to EKS
   ```bash
   kubectl apply -f k8s/
   ```

9. Get External IP
   ```bash
   kubectl get services -n ingress-nginx
   ```

### Cleanup
- Delete EKS Cluster
  ```bash
  eksctl delete cluster --name go-microservices-cluster --region us-east-1
  ```

## Endpoints
- Auth Service: `/auth/*`
- API Service: `/api/*`

## Technologies
- Go
- Kubernetes
- Docker
- PostgreSQL
- NGINX Ingress Controller
- AWS EKS

## Contributing
1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request
