# User System

## Project Overview

This project is a user management system based on Spring Boot 3, supporting user registration, login, information management, role-based access control, and data export. The backend uses PostgreSQL as the main database and Redis for caching and distributed scenarios. JWT is used for stateless authentication and role-based authorization.

## Tech Stack

* Java 17
* Spring Boot 3.5.4
* MyBatis
* PostgreSQL 17.x
* MySQL 8.0.x
* Redis
* Spring Security
* EasyExcel (Alibaba Excel export)
* Lombok
* JJWT (JWT token)
* Maven 3.9.x

## Project Structure

```
backend/
  src/
    main/
      java/com/bryan/system/
        config/         # Configuration classes (security, Redis, MyBatis, etc.)
        controller/     # RESTful controllers
        domain/         # Entities, request/response objects, VO
        filter/         # JWT authentication filter
        handler/        # Global exception handler
        mapper/         # MyBatis mapper interfaces
        service/        # Service layer
        util/           # Utility classes (JWT, HTTP, etc.)
      resources/
        application.yaml
        application-dev.yaml
        application-mysql.yaml
        mapper/         # MyBatis mapper xmls
    test/
      java/com/bryan/system/
        UserSystemApplicationTests.java
  pom.xml
  mvnw
frontend/
  src/
  package.json
```

## Requirements

* JDK 17+
* Maven 3.9.9+
* PostgreSQL 17.x/MySQL 8.0.x
* Redis 6.x or above

## Configuration

* Update database and Redis settings in `backend/src/main/resources/application-dev.yaml`.
* General settings (logging, MyBatis, etc.) are in `backend/src/main/resources/application.yaml`.
* Database schema scripts are in [`sql/create_table.sql`](sql/create_table.sql).

## Getting Started

> Using PostgreSQL as example

1. Initialize the PostgreSQL database by running the schema script:

   ```sh
   psql -U postgres -d postgres -f sql/create_table.sql
   ```
2. Start the Redis service.
3. Build and run the project with Maven:

   ```sh
   cd backend
   ./mvnw spring-boot:run
   ```

   Or run the packaged jar:

   ```sh
   cd backend
   ./mvnw clean package
   java -jar target/user-system-0.0.1-SNAPSHOT.jar
   ```

## Main APIs

* User registration: `POST /api/auth/register`
* User login: `POST /api/auth/login`
* Get all users: `GET /api/user/all` (admin only)
* Search users: `POST /api/user/search`
* User update, role change, password update, ban/unban, logical delete, etc. are detailed in [`UserController`](backend/src/main/java/com/bryan/system/controller/user/UserController.java)
* Export user data: `GET /api/user/export/all`, `POST /api/user/export/field` (admin only)

## Notes

* For production, inject JWT secret via configuration file instead of hardcoding.
* Global exception handling is in [`GlobalExceptionHandler`](backend/src/main/java/com/bryan/system/handler/GlobalExceptionHandler.java).
* Logical delete field is `deleted`: 0 means active, 1 means deleted.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
