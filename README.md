# ziyadbook

Go REST API with Gin, clean architecture, and TDD-first design. Runs in Docker Compose with Nginx, MySQL, and Redis.

## Architecture

- **Nginx** (`:8080`) → reverse proxy to Go API
- **Go API** (`:8080` in containers)
  - handler → service → repository layers
  - Gin router, `/health`, and `users` CRUD
- **MySQL 8.4** (`:3306`) with migrations
- **Redis 7.4** (`:6379`)

## Quick start

### Prerequisites

- Docker and Docker Compose
- Make (optional, for convenience)

### 1) Clone and prepare

```bash
git clone <repo>
cd ziyadbook
cp .env.example .env
# Edit .env if needed (defaults work for Docker)
```

### 2) Production-like mode

```bash
make up
```

- Nginx proxies to `api` (runtime binary)
- Access: `http://localhost:8080/health`

### 3) Development mode (hot reload)

```bash
make dev
```

- Nginx proxies to `api-dev`
- Air watches `*.go` and rebuilds/restarts automatically
- Access: `http://localhost:8080/health`

### 4) Run tests (inside container)

```bash
make tidy   # generate go.sum
make test   # run unit tests
```

### 5) Stop

```bash
make down
```

## Makefile commands

| Command | What it does |
|---------|--------------|
| `make up` | Starts prod-like stack (`nginx`, `api`, `mysql`, `redis`) |
| `make dev` | Starts dev stack with hot reload (`nginx`, `api-dev`, `mysql`, `redis`) |
| `make down` | Stops and removes all containers/volumes |
| `make test` | Runs Go unit tests inside `api-dev` |
| `make tidy` | Runs `go mod tidy` inside `api-dev` (generates `go.sum`) |

## API endpoints

- `GET /health` – health check
- `POST /users` – create a user (`{"email":"…","name":"…"}`)
- `GET /users/:id` – fetch a user by ID
- `GET /users?limit=N` – list users (default limit 20)
- `POST /borrow` – borrow a book (`{"book_id":…,"member_id":…}`) with atomic stock/quota validation

## Development notes

- **Hot reload**: Air watches `*.go` in `api-dev`; changes rebuild/restart automatically
- **MySQL healthcheck**: `api`/`api-dev` wait for MySQL to be healthy before starting
- **TDD**: Unit tests exist in `*_test.go` files; use `make test` to run them
- **Profiles**: `api` is `prod` profile; `api-dev` is `dev` profile
