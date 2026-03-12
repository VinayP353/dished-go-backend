# Dished Go Backend

A modular Go backend application with PostgreSQL and Redis, fully containerized with Docker.

## Features

- Clean architecture with separation of concerns
- PostgreSQL database with GORM
- Redis caching layer
- RESTful API with Gin framework
- Docker containerization
- Health check endpoints
- CORS middleware
- Request logging

## Project Structure

```
dished-go-backend/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── database/             # Database connections
│   ├── handlers/             # HTTP handlers
│   ├── middleware/           # HTTP middleware
│   ├── models/               # Data models
│   ├── repository/           # Data access layer
│   └── service/              # Business logic layer
├── docker-compose.yml
├── Dockerfile
└── go.mod
```

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)

### Running with Docker

```bash
docker-compose up --build
```

The application will be available at `http://localhost:8080`

### API Endpoints

- `GET /health` - Health check
- `POST /api/v1/users` - Create user
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Example Request

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","name":"John Doe"}'
```

## Development

### Local Development

```bash
go mod download
go run ./cmd/api/main.go
```

### Stop Services

```bash
docker-compose down
```

### Clean Up

```bash
make clean
```
