# 🖥️ Backend — Face Recognition Attendance System

REST API backend service built with **Go (Fiber)** for a face recognition-based attendance system. Handles user management, attendance tracking, Jira integration for KPI monitoring, and communicates with the ML service for face inference.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.25 |
| Framework | [Fiber v2](https://github.com/gofiber/fiber) |
| Database | PostgreSQL 17 |
| Object Storage | MinIO |
| Cache | Redis 7 |
| Auth | JWT (golang-jwt) |
| Scheduler | gocron |
| Documentation | Swagger (swaggo) |
| Containerization | Docker |
| Project Management | Jira REST API |

## Architecture

The project follows **Clean Architecture** with clear separation of concerns:

```
backend-facerecognition/
├── app/                    # Application entry point
│   └── main.go             # Server bootstrap, DI, scheduler
├── facerpca/               # Core feature module
│   ├── delivery/           # HTTP handlers (Fiber routes)
│   │   └── http/
│   ├── repository/         # Data access layer
│   │   ├── sql/            # PostgreSQL queries
│   │   ├── minio/          # Object storage operations
│   │   ├── redis/          # Cache operations
│   │   └── jira/           # Jira API integration
│   └── usecase/            # Business logic
├── domain/                 # Domain models & interfaces
├── constant/               # Application constants
├── helper/                 # Utility functions & error handling
├── scripts/                # Python utilities for dataset prep
│   ├── register_face.py
│   ├── video_to_dataset.py
│   └── video_manual_extraction.py
├── deploy/                 # Deployment configuration
│   ├── docker-compose.yml  # Full stack (Postgres, Redis, MinIO, Backend)
│   ├── nginx.conf          # Reverse proxy config
│   └── deploy.sh           # Deployment script
├── Dockerfile.backend      # Multi-stage Go build
├── go.mod / go.sum         # Go dependencies
└── .env.example            # Environment variable template
```

## Features

- **Face Registration & Inference** — communicates with ML service for face recognition
- **Attendance Management** — clock-in/clock-out with face verification
- **Automated Scheduling**:
  - `05:00` — Initialize daily attendance status
  - `07:00` — Overnight Jira catch-up for overtime workers
  - `09:00` — Mark absent (morning cutoff)
  - `Every 30min` — Poll Jira tasks for active users
  - `20:00` — Force clock-out & KPI accumulation
- **Jira Integration** — automatic task polling and KPI tracking
- **Dashboard & Reports** — attendance statistics and Excel export
- **Redis Caching** — session management and performance optimization

## Getting Started

### Prerequisites
- Go 1.25+
- Docker & Docker Compose
- PostgreSQL 17, Redis 7, MinIO (or use Docker Compose)

### Configuration

1. Copy the example config file:
   ```bash
   cp config.yaml.example config.yaml
   cp .env.example .env
   ```

2. Fill in your actual credentials in both files.

### Run with Docker (Recommended)

```bash
# Create the shared Docker network (first time only)
docker network create app_network

# Start the full stack
cd deploy
docker compose up -d --build
```

This starts PostgreSQL, Redis, MinIO, and the backend service.

### Run Locally

```bash
# Install dependencies
go mod download

# Run the server
go run app/main.go -c config.yaml
```

The server will start on port `5000` by default.

## API Documentation

Swagger docs are available at `/swagger/` when the server is running.

## Environment Variables

See [`.env.example`](.env.example) for the full list of required environment variables:
- `DATABASE_*` — PostgreSQL connection
- `MINIO_*` — Object storage credentials
- `REDIS_*` — Cache configuration
- `JIRA_*` — Jira API integration

## Related Services

- [**Frontend**](https://github.com/cenjaa/frontend-facerecognition) — React-based kiosk UI
- [**ML Service**](https://github.com/cenjaa/ml-service) — RPCA + PCA + SVM face recognition pipeline
