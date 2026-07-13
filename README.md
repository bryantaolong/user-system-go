# User System

## Project Overview

This project is a user management system based on Go (Gin framework), supporting user registration, login, information management, role-based access control, and data export. The backend uses PostgreSQL as the main database and Redis for caching and distributed scenarios. JWT is used for stateless authentication and role-based authorization.

## Tech Stack

* Go 1.22
* Gin - Web framework
* GORM - ORM
* go-redis/v9 - Redis client
* golang-jwt/v5 - JWT authentication
* Viper - Configuration management
* Zap - Logging
* excelize - Excel export
* bcrypt - Password encryption

## Project Structure

```
backend/
  cmd/server/           # Application entry point
  internal/
    config/            # Configuration loading
    handler/           # HTTP handlers
    middleware/        # Middleware (auth, CORS, error handling)
    model/             # Data models & DTOs
    pkg/               # Utility packages (JWT, Redis, HTTP, response)
    repository/        # Data access layer
    service/           # Business logic layer
  config.yaml          # Configuration file
  go.mod
frontend/
  src/
  package.json
```

## Requirements

* Go 1.22+
* PostgreSQL 17.x / MySQL 8.0.x
* Redis 6.x or above

## Configuration

* Edit `backend/config.yaml` to configure database, Redis, JWT, and other parameters.
* Database schema scripts are in [`sql/create_table.sql`](sql/create_table.sql).

## Getting Started

> Using PostgreSQL as example

1. Initialize the PostgreSQL database by running the schema script:

   ```sh
   psql -U postgres -d postgres -f sql/create_table.sql
   ```
2. Start the Redis service.
3. Install dependencies and run the project:

   ```sh
   cd backend
   go mod tidy
   go run cmd/server/main.go
   ```

## Main APIs

* User registration: `POST /api/auth/register`
* User login: `POST /api/auth/login`
* Get current user: `GET /api/auth/me`
* Get all users: `GET /api/users` (admin only)
* Search users: `POST /api/users/search`
* User update, role change, password update, ban/unban, logical delete, etc.
* Export user data: `GET /api/users/export` (admin only)
* User profiles: avatar upload, profile CRUD

## Notes

* For production, inject JWT secret via configuration file instead of hardcoding.
* Logical delete field is `deleted`: 0 means active, 1 means deleted.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
