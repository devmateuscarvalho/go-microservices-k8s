# Go Microservices with JWT Authentication

Este projeto consiste em dois microserviços em Go:
1. Auth Service: Serviço de autenticação com JWT
2. API Service: API protegida que utiliza o JWT do Auth Service

## Requisitos

- Go 1.21+
- Docker
- Kubernetes (Minikube ou similar)
- PostgreSQL

## Estrutura do Projeto

```
.
├── auth-service/         # Serviço de autenticação
├── api-service/         # API protegida
├── k8s/                # Configurações Kubernetes
└── docker-compose.yml  # Configuração para desenvolvimento local
```

## Como Executar

1. Iniciar os serviços localmente:
```bash
docker-compose up
```

2. Aplicar as configurações Kubernetes:
```bash
kubectl apply -f k8s/
```
