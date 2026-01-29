# User Restriction Manager

A robust Go service for managing user account restrictions, appeals, and audit trails.

## Features

- **Restrictions Management**: Apply warnings, temporary bans, and permanent bans.
- **Appeals Workflow**: Users can submit appeals, and admins can review/approve/reject them.
- **Audit Logging**: (Planned) Track all administrative actions.
- **High Performance**: Uses Redis for caching active restriction checks.
- **Clean Architecture**: Domain-driven design with hexagonal architecture.

## Tech Stack

- **Go 1.23**
- **PostgreSQL 16** (Primary Data Store)
- **Redis 7** (Caching)
- **Docker & Docker Compose**
- **Chi Router**
- **Go-Migrate** (planned integration)

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.23+ (for local development)

### Running the Application

1. **Start Infrastructure**:
   ```bash
   docker-compose up -d
   ```

2. **Run the Server**:
   ```bash
   go run cmd/server/main.go
   ```

   The server will start on port `8080`.

### API Endpoints

#### Restrictions

- `POST /v1/restrictions` - Apply a new restriction
- `GET /v1/restrictions/{userID}` - Get active restrictions for a user
- `DELETE /v1/restrictions/{id}` - Revoke a restriction

#### Appeals

- `POST /v1/appeals` - Submit an appeal for a restriction
- `POST /v1/appeals/{id}/review` - Review an appeal (Approve/Reject)

## Environment Variables

See `.env` for configuration options.
