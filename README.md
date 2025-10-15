# Contact Management API

A RESTful API for contact management built with Go, Gin, and MySQL.

## Features

- User registration and authentication with JWT
- Contact CRUD operations
- Favorite contacts
- Search and pagination
- MySQL database with proper indexing
- Docker support

## Prerequisites

- Go 1.24+
- MySQL 8.0+
- Docker (optional)

## Database Setup

### Using Migrations (Recommended)

The project now uses a proper migration system for database schema management:

```bash
# Apply all migrations
make migrate-up

# Check migration status
make migrate-status

# Rollback if needed
make migrate-down
```

See [`MIGRATIONS.md`](MIGRATIONS.md) for detailed migration documentation.

### Legacy Setup (setup_db.sh - No Longer Used)

The `setup_db.sh` script is **deprecated** and should not be used. It has been replaced by the migration system above.

## Environment Configuration

Copy `.env.example` to `.env` and update the values:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_mysql_user
DB_PASSWORD=your_mysql_password
DB_NAME=getcontact

# JWT Configuration
JWT_SECRET=your_jwt_secret_key

# Server Configuration
PORT=8080
ENVIRONMENT=development
```

## Installation & Running

1. Install dependencies:

   ```bash
   go mod tidy
   ```

2. Run the application:

   ```bash
   go run ./cmd/server
   ```

3. The API will be available at `http://localhost:8080`

## API Endpoints

### Authentication

- `POST /api/v1/register` - User registration
- `POST /api/v1/login` - User login

### Contacts (Protected routes)

- `GET /api/v1/contacts` - List contacts with search/pagination
- `POST /api/v1/contacts` - Create new contact
- `GET /api/v1/contacts/:id` - Get contact details
- `PUT /api/v1/contacts/:id` - Update contact
- `DELETE /api/v1/contacts/:id` - Delete contact

### User Profile

- `GET /api/v1/profile` - Get user profile
- `PUT /api/v1/profile` - Update user profile

## Database Schema

### Users Table

- `id` (Primary Key, Auto Increment)
- `full_name` (Indexed)
- `email` (Unique, Indexed)
- `phone` (Indexed)
- `password` (Hashed)
- `avatar_url`
- `created_at` (Indexed)
- `updated_at`

### Contacts Table

- `id` (Primary Key, Auto Increment)
- `user_id` (Foreign Key to users.id, Indexed)
- `full_name` (Indexed)
- `phone` (Indexed)
- `email` (Indexed)
- `favorite` (Indexed)
- `created_at` (Indexed)
- `updated_at`

### Indexes

- Single column indexes on frequently queried fields
- Composite indexes for user-specific queries
- Foreign key constraints with CASCADE delete

## Docker Support

Build and run with Docker:

```bash
docker-compose up --build
```

## API Documentation

Complete API documentation with cURL examples is available in [`api_examples.md`](api_examples.md).

### Quick API Test

Run the included test script to verify all endpoints:

```bash
./test_api.sh
```

This script will:

- Test the health endpoint
- Register a test user
- Login and get JWT token
- Test all protected endpoints (profile, contacts CRUD)

**Note:** Make sure the server is running and migrations are applied before running the test script.

## Project Structure

```
├── cmd/server/          # Application entry point
├── configs/             # Configuration management
├── internal/
│   ├── app/
│   │   ├── handlers/    # HTTP handlers
│   │   ├── models/      # Data models
│   │   ├── repository/  # Database layer
│   │   ├── routes/      # Route definitions
│   │   └── service/     # Business logic
│   ├── logger/          # Logging middleware
│   └── middleware/      # HTTP middleware
├── pkg/
│   ├── db/             # Database utilities
│   └── redis/          # Redis client
├── database_schema.sql # MySQL schema
├── setup_db.sh        # Database setup script
└── docker-compose.yml # Docker configuration
```
