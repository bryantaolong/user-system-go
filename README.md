# User System (Go Version)

[中文说明请见这里 (Chinese README here)](./README_zh.md)

## Project Overview

This project is a user management system based on the Go programming language. It supports user registration, login, information management, role-based access control, and data export. The backend uses PostgreSQL as the main database and Redis for caching. JWT is used for stateless authentication and role-based authorization.

## Tech Stack

- Go 1.20+
- Gin Web Framework
- GORM (ORM framework)
- PostgreSQL 17.x
- Redis 6.x or above
- JWT (JSON Web Token)

## Project Structure

```
cmd/
  main.go                # Application entry point
internal/
  config/                # Configuration
  handler/               # Route handlers (user, auth, etc.)
  middleware/            # Middleware (auth, permission, etc.)
  model/
    entity/              # Entity definitions
    request/             # Request structs
    response/            # Response structs
  router/                # Route registration
  service/               # Business logic
    redis/               # Redis-related services
pkg/
  db/                    # Database initialization
  http/                  # HTTP utilities
  jwt/                   # JWT utilities
go.mod, go.sum           # Go dependency management
README.md, README_zh.md  # Project documentation
LICENSE                  # License
```

## Main Features

- User registration, login, logout
- User information query and modification
- Password change
- Role management and access control
- User disable/enable, logical deletion
- User data pagination and search
- User data export (e.g., CSV/Excel)
- JWT authentication and middleware
- Redis cache support

## Requirements

- Go 1.20 or above
- PostgreSQL 17.x
- Redis 6.x or above

## Configuration

- Database, Redis, and other configurations should be set in `internal/config/config.go` or via environment variables.
- Other general configurations can be found in the `internal/config` directory.

## Getting Started

1. Initialize the database (PostgreSQL). Table creation SQL can be auto-generated from `model/entity` structs or written manually.
2. Start the Redis service.
3. Install dependencies and run:

   ```sh
   go mod tidy
   go run cmd/main.go
   ```

   Or build and run:

   ```sh
   go build -o user-system-go cmd/main.go
   ./user-system-go
   ```

## Main APIs

- User registration: `POST /api/auth/register`
- User login: `POST /api/auth/login`
- Get all users: `GET /api/user/all` (admin only)
- Search users: `POST /api/user/search`
- User update, role change, password update, ban/unban, logical delete, etc. are detailed in `internal/handler/user_handler.go`
- Export user data: e.g., `GET /api/user/export/all`, `POST /api/user/export/field` (admin only)

## Notes

- It is recommended to inject the JWT secret via configuration file or environment variable instead of hardcoding.
- Global exception handling and unified response format can be implemented in `internal/handler` or middleware.
- Logical delete field is recommended as `deleted`: 0 means active, 1 means deleted.

## License

This project is licensed under the MIT License.
See [LICENSE](LICENSE) for