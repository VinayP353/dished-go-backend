# Dished Go Backend

A modular Go backend for the Dished chef platform, built with PostgreSQL, Redis, and Docker.

## Features

- Clean architecture (handler → service → repository)
- PostgreSQL with GORM (auto-migration)
- Redis caching layer
- RESTful API with Gin
- Swagger UI documentation
- Docker Compose setup
- CORS middleware + request logging

## Project Structure

```
dished-go-backend/
├── cmd/api/main.go           # Entry point + Swagger general info
├── docs/                     # Auto-generated Swagger docs
├── internal/
│   ├── config/               # Env config
│   ├── database/             # Postgres + Redis connections
│   ├── handlers/             # HTTP handlers
│   ├── middleware/           # CORS, logger
│   ├── models/               # GORM models + DTOs
│   ├── repository/           # Data access layer
│   └── service/              # Business logic
├── docker-compose.yml
├── Dockerfile
└── go.mod
```

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.23+ (for local development)
- `swag` CLI (for regenerating docs)

### Run with Docker

```bash
docker compose up --build
```

App runs at `http://localhost:8081`

### Swagger UI

```
http://localhost:8081/swagger/index.html
```

### Regenerate Swagger Docs

After changing handler annotations or models:

```bash
swag init -g cmd/api/main.go --output docs
```

---

## API Endpoints

### Health

| Method | Path      | Description  |
|--------|-----------|--------------|
| GET    | /health   | Health check |

---

### Auth

| Method | Path                  | Description       | Body                                      |
|--------|-----------------------|-------------------|-------------------------------------------|
| POST   | /api/v1/auth/register | Register a chef   | `username`, `password`, `email`, `first_name`, `last_name` |
| POST   | /api/v1/auth/login    | Login as a chef   | `username`, `password`                    |

**Register example:**
```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "chefjohn",
    "password": "SecureP@ss123",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

**Login example:**
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "chefjohn", "password": "SecureP@ss123"}'
```

---

### Chefs

| Method | Path                       | Description              |
|--------|----------------------------|--------------------------|
| GET    | /api/v1/chefs              | Get all chefs            |
| GET    | /api/v1/chefs/:id          | Get chef by ID           |
| GET    | /api/v1/chefs/usernames    | Get all usernames        |
| PUT    | /api/v1/chefs/:id          | Update chef              |
| PUT    | /api/v1/chefs/:id/profile  | Update chef's profile    |
| DELETE | /api/v1/chefs/:id          | Delete chef              |

**Update chef example:**
```bash
curl -X PUT http://localhost:8081/api/v1/chefs/1 \
  -H "Content-Type: application/json" \
  -d '{"email": "newemail@example.com", "status": "active"}'
```

**Update profile example:**
```bash
curl -X PUT http://localhost:8081/api/v1/chefs/1/profile \
  -H "Content-Type: application/json" \
  -d '{
    "preferred_name": "Chef John",
    "address": "123 Culinary Lane",
    "description": "Passionate about Italian cuisine"
  }'
```

---

### Chef Profiles

Standalone profile management (used internally or for admin purposes).

| Method | Path                        | Description              |
|--------|-----------------------------|--------------------------|
| POST   | /api/v1/chef-profiles       | Create a profile         |
| GET    | /api/v1/chef-profiles       | Get all profiles         |
| GET    | /api/v1/chef-profiles/:id   | Get profile by ID        |
| PUT    | /api/v1/chef-profiles/:id   | Update profile           |
| DELETE | /api/v1/chef-profiles/:id   | Delete profile           |

---

### Users

| Method | Path                  | Description       |
|--------|-----------------------|-------------------|
| POST   | /api/v1/users         | Create user       |
| GET    | /api/v1/users         | Get all users     |
| GET    | /api/v1/users/:id     | Get user by ID    |
| PUT    | /api/v1/users/:id     | Update user       |
| DELETE | /api/v1/users/:id     | Delete user       |

---

## Password Requirements

Passwords must:
- Be at least 8 characters
- Contain at least 6 letters
- Contain at least 1 number
- Contain at least 1 special character (`!@#$%^&*` etc.)

---

## Development

```bash
# Download dependencies
go mod download

# Run locally (requires Postgres + Redis running)
go run ./cmd/api/main.go

# Stop Docker services
docker compose down

# Clean up volumes
make clean
```
