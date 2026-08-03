# User System

## Project Overview

This project is a user management system based on Go (Gin framework), supporting user registration, login, information management, role-based access control, and data export. The backend uses PostgreSQL as the main database and Redis for caching and distributed scenarios. JWT is used for stateless authentication and role-based authorization.

## Tech Stack

### Backend

* Go 1.25
* Gin - Web framework
* `database/sql` + GORM - ORM
* go-redis/v9 - Redis client
* golang-jwt/v5 - JWT authentication
* Viper - Configuration management
* Zap - Logging
* excelize - Excel export
* bcrypt - Password encryption

### Frontend

* Vue 3
* TypeScript
* Vite
* Arco Design Vue
* Pinia
* Pinia Plugin Persist
* Vue Router
* Axios

## Project Structure

```
backend/
  main.go              # Application entry point
  auth/                # Auth module (handler, service, repository, middleware)
  cache/               # Redis client
  config/              # Configuration loading
  config.yaml          # Configuration file
  middleware/          # Global middleware (CORS, error handling)
  model/               # Data models & DTOs
  pkg/                 # Utility packages (JWT, Redis, HTTP, response)
  response/            # Response helpers
  system/              # System module (logs, etc.)
  user/                # User module (handler, service, repository, file)
  go.mod

frontend/
  src/
    api/               # API request modules
    assets/            # Static assets
    components/        # Shared components
    router/            # Route definitions & guards
    stores/            # Pinia state stores
    styles/            # Global styles
    types/             # TypeScript type definitions
    utils/             # Utility helpers
    views/             # Page views
      login/           # Login page
      register/        # Registration page
      profile/         # User profile pages
      admin/           # Admin-only pages
        users/         # User management
        logs/          # System logs
  package.json
  vite.config.ts
  tsconfig.json
```

## Requirements

* Go 1.25+
* Node.js 18+
* PostgreSQL 17.x / MySQL 8.0.x
* Redis 6.x or above

## Configuration

* Edit `backend/config.yaml` to configure database, Redis, JWT, server port, and other parameters.
* Database schema scripts are in [`sql/create_table.sql`](sql/create_table.sql).

Example `config.yaml`:

```yaml
server:
  port: 8080

database:
  driver: postgres
  host: localhost
  port: 5432
  username: postgres
  password: "postgres"
  dbname: user_system
  sslmode: disable

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret_key: "your-secret-key"
  expiration_ms: 86400000
  token_prefix: "Bearer "

cors:
  allowed_origins: "http://localhost:5173"

logging:
  level: info
  file: logs/user-system.log
```

## Getting Started

> Using PostgreSQL as example

1. Initialize the PostgreSQL database by running the schema script:

   ```sh
   psql -U postgres -d postgres -f sql/create_table.sql
   ```

2. Start the Redis service.
3. Start the backend service:

   ```sh
   cd backend
   go mod tidy
   go run main.go
   ```

4. Start the frontend service in a separate terminal:

   ```sh
   cd frontend
   npm install
   npm run dev
   ```

5. Visit `http://localhost:5173` (frontend defaults to Vite dev server).

### Build

```sh
# Frontend
cd frontend && npm run build

# Backend
cd backend && go build -o bin/server main.go
```

## Main APIs

* User registration: `POST /api/auth/register`
* User login: `POST /api/auth/login`
* Validate token: `GET /api/auth/validate`
* Get current user: `GET /api/auth/me`
* Change password: `PUT /api/auth/password`
* Delete account: `DELETE /api/auth`
* Logout: `GET /api/auth/logout`

* Get all users: `GET /api/users` (admin only)
* Get user by ID: `GET /api/users/:userId` (admin only)
* Get user by username: `GET /api/users/username/:username` (admin only)
* Search users: `POST /api/users/search` (admin only)
* Create user: `POST /api/users` (admin only)
* Update user: `PUT /api/users/:userId`
* Change user roles: `PUT /api/users/roles/:userId` (admin only)
* Reset password: `PUT /api/users/password/:userId` (admin only)
* Block user: `PUT /api/users/block/:userId` (admin only)
* Unblock user: `PUT /api/users/unblock/:userId` (admin only)
* Delete user: `DELETE /api/users/:userId` (admin only, logical delete)
* Export users: `GET /api/users/export` (admin only)

* Upload avatar: `POST /api/user-profiles/avatar`
* Get profile by user ID: `GET /api/user-profiles/:userId`
* Get profile by real name: `GET /api/user-profiles/name/:realName`
* Get current user profile: `GET /api/user-profiles/me`
* Update profile: `PUT /api/user-profiles`

* List roles: `GET /api/user-roles` (admin only)

* List latest logs: `GET /api/admin/logs` (admin only)
* List log files: `GET /api/admin/logs/files` (admin only)

## Frontend

* Vue 3 + TypeScript SPA with Vite and Arco Design Vue.
* Authentication flows: login, registration, logout, and protected route guards.
* User profile management: basic info, security settings, login history.
* Admin panel: user management with search, update, role change, ban/unban, delete; system log viewing.
* Axios interceptor for JWT injection and unified error handling.

## Notes

* For production, inject JWT secret via configuration file instead of hardcoding.
* Logical delete field is `deleted`: 0 means active, 1 means deleted.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
