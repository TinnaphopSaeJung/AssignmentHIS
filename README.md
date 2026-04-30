# HIS Backend System

A backend service for a Hospital Information System (HIS) built with Go, Gin, PostgreSQL, and Docker.

This system manages patients, hospitals, and staff, including multi-hospital patient mapping.

---

## Tech Stack

- Go (Golang)
- Gin Web Framework
- PostgreSQL
- Docker & Docker Compose
- Nginx (Reverse Proxy)

---

## Project Structure
`
├── database/
│   └── init.sql            # SQL script for initializing the database (create tables, seed data)
│
├── internal/
│   ├── clients/            # External service clients (e.g., Hospital A API)
│   ├── config/             # Environment variables and application configuration
│   ├── database/           # Database connection setup (PostgreSQL)
│   ├── handler/            # HTTP request/response handlers
│   ├── middleware/         # Middleware (e.g., JWT authentication)
│   ├── models/             # Structs mapping to database tables
│   ├── repository/         # Database query layer
│   ├── service/            # Business logic layer
│   └── nginx/              # Nginx configuration
│
├── pkg/
│   └── utils/              # Utility functions (e.g., response, validation, JWT)
│
├── tests/                  # Unit tests
│
├── docker-compose.yml      # Multi-container orchestration
├── Dockerfile              # Build configuration for Go service image
├── .env                    # Environment variables for local development
├── .env.docker             # Environment variables for Docker environment
├── go.mod                  # Go module dependencies
└── go.sum                  # Dependency checksums
`

---

## Getting Started

### 1. Clone the repository
`git clone https://github.com/TinnaphopSaeJung/AssignmentHIS.git`
`cd AssignmentHIS`
### 2. Run with Docker
`docker compose up --build`

---

## API Endpoints
### Staff
- POST /staff/create
- POST /staff/login
### Patient
- POST /patient/search
- POST /patient/search-from-external

---

## Database
- PostgreSQL is used as the main database
- Schema is initialized via `database/init.sql`
- Core tables:
  - patients
  - hospitals
  - staffs
  - patient_hospitals_mapping

---

 ## Architecture
- Nginx acts as a reverse proxy
- Go (Gin) handles API logic
- PostgreSQL handles data persistence
- Docker orchestrates all services

---

## Documentation
- Development Planning Docs (Google Docs): https://docs.google.com/document/d/1HkA6Fx-OmZf4JzSPBcW6mwQjyl145fc7sgIb1v5jfDk/edit?usp=sharing

------------------------------------------------------------------------------------------------------------------------------------------------
## Author
- Tinnaphop Sae-jung (Cake)
- sj.tinnaphop@gmail.com
