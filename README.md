# Job Scheduler

[![CI](https://github.com/ckng0221/job-scheduler/actions/workflows/ci.yml/badge.svg)](https://github.com/ckng0221/job-scheduler/actions/workflows/ci.yml)

Job Scheduler is a proof of concept (POC) distributed job scheduler written in [Go](https://go.dev/) and [TypeScript](https://www.typescriptlang.org/). The application allows users to submit their jobs and run on desired schedule, for both one-time trigger and recurrence schedule. It consists of four main components:

- `API`: Responsible for creating jobs and executions, user authentication, and other business logic.
- `Scheduler`: Polls the latest active jobs from the database and submits them to the job queue.
- `Worker`: Receives job events from the job queue and executes each job.
- `UI`: Provides a user interface for logging in and creating scheduled jobs.

## Tech Stacks

### API

- Programming Language: [Go](https://go.dev/)
- Server Framework: [Gin](https://pkg.go.dev/github.com/gin-gonic/gin)
- ORM: [Gorm](https://gorm.io/)
- Database: [MySQL](https://www.mysql.com/)
- Authentication Protocol: [OIDC](https://openid.net/developers/how-connect-works/)
- Identity Provider: [Google](https://developers.google.com/identity)

### Scheduler, Worker

- Programming Language: [Go](https://go.dev/)
- Message Broker: [RabbitMQ](https://www.rabbitmq.com/)

### UI

- Programming Language: [TypeScript](https://www.typescriptlang.org/)
- Web Framework: [Next.js](https://nextjs.org/)
- CSS Framework: [Tailwind CSS](https://tailwindcss.com/)
- UI Library: [Material UI](https://mui.com/)

### Build

- CI Platform: [GitHub Actions](https://github.com/features/actions)
- Build System: [Turborepo](https://turbo.build/)
- Multi-container Tool: [Docker Compose](https://docs.docker.com/compose/)

## Getting Started

### Installation

```bash
# At the project root

$ npm install

# Install Go dev dependencies
$ go get github.com/githubnemo/CompileDaemon
$ go install github.com/githubnemo/CompileDaemon
```

Before running application, rename the `.env.example` files to `.env`, and update the environment variables accordingly.

## Run application

### On local

To run the application locally, ensure that `MySQL` and `RabbitMQ` are installed beforehand.

```bash
# At the project root
# This will run all four modules simultaneously.

# Development mode
npm run dev

# Build
npm run build

# Production mode
npm run start

# Alternatively, you can navigate to the root of each application (e.g., ./apps/api) and run the npm scripts to run the particular application only.
```

### With docker and docker compose

To run the application using Docker, ensure that `Docker` and `Docker Compose` are installed beforehand.

```bash
# Create docker images and run docker containers in detached mode
docker compose up -d

# Stop and remove containers
docker compose down
```

### Run with Kubernetes

```bash
# update configmap and rename api.configmap.example.yaml to api.configmap.yaml
# apply yaml files in k8s folder
kubectl apply -f k8s

kubectl port-forward services/api 8000:8000
kubectl port-forward services/ui 3000:3000

# TO access from minikube
minikube service api --url
minikube service ui --url
```

# Contributing

Contributions are welcome! If you'd like to contribute to this project, please follow the guidelines outlined in CONTRIBUTING.md.

# License

This project is licensed under the MIT License.
