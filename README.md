# Go Fiber GORM Boilerplate

A production-ready RESTful API boilerplate using Go Fiber and GORM with advanced features for large-scale applications.

## Features

- **Go Fiber**: High-performance web framework
- **GORM**: Feature-rich ORM with PostgreSQL support
- **Clean Architecture**: Following best practices of separation of concerns
- **API Versioning**: Support for multiple API versions
- **Authentication**: JWT-based authentication system
- **Authorization**: Role-based access control
- **Middleware**: Authentication, logging, error handling, rate limiting, validation
- **Configuration**: Environment-specific configuration using .env files
- **Migrations**: Database migration system with rollback support
- **Error Handling**: Centralized error handling with detailed error types
- **Logging**: Structured logging with different levels
- **Caching**: Redis-based caching for improved performance
- **Rate Limiting**: Protection against abuse and DoS attacks
- **Request Validation**: Automatic request validation using struct tags
- **Metrics & Monitoring**: Prometheus metrics and monitoring dashboard
- **Worker Pool**: Background task processing system
- **Health Checks**: Advanced health check system for monitoring
- **Containerization**: Docker and Docker Compose support
- **CI/CD**: GitHub Actions workflow for continuous integration and deployment
- **Swagger Documentation**: API documentation with Swagger
- **Testing**: Comprehensive testing utilities for unit and integration tests
- **Transaction Management**: Database transaction handling

## Project Structure

```
├── .github/                      # GitHub Actions workflows
├── cmd/                          # Application entry points
│   └── main.go                   # Main application
├── config/                       # Configuration
│   ├── config.go                 # Configuration handling
│   └── env_loader.go             # Environment configuration
├── internal/                     # Private application code
│   ├── handler/                  # HTTP handlers
│   ├── middleware/               # Custom middleware
│   ├── model/                    # Data models
│   ├── repository/               # Database operations
│   └── service/                  # Business logic
├── migrations/                   # Database migrations
├── pkg/                          # Public libraries
│   ├── auth/                     # Authentication utilities
│   ├── cache/                    # Caching utilities
│   ├── docs/                     # API documentation
│   ├── errors/                   # Error handling
│   ├── logger/                   # Logging utilities
│   ├── monitoring/               # Metrics and monitoring
│   └── worker/                   # Background processing
├── routes/                       # Route definitions
├── test/                         # Testing utilities
├── .env.example                  # Example environment variables
├── docker-compose.yml            # Docker Compose configuration
├── Dockerfile                    # Docker image definition
├── go.mod                        # Go module file
├── go.sum                        # Go module checksums
├── Makefile                      # Makefile for common operations
└── README.md                     # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Redis (optional but recommended)
- Docker and Docker Compose (for containerized development)

### Installation

#### Local Development

1. Clone the repository
2. Configure your `.env` file (use `.env.example` as a template)
```bash
cp .env.example .env
```
3. Install dependencies
```bash
go mod tidy
```
4. Start the application
```bash
go run cmd/main.go
```

#### Using Docker

```bash
# Start all services with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

### Using Make Commands

This project includes a Makefile for common operations:

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Run linter
make lint

# Deploy with Docker Compose
make docker-compose-up
```

## API Documentation

API documentation is available via Swagger at:
```
http://localhost:8080/swagger/
```

## Monitoring

Monitoring endpoints are available at:
- Dashboard: `http://localhost:8080/dashboard`
- Metrics: `http://localhost:8080/metrics` (Prometheus compatible)

## Health Checks

Health check endpoints:
- Basic: `http://localhost:8080/health`
- Detailed: `http://localhost:8080/health/details`

## Authentication and Authorization

The API uses JWT for authentication. Protected routes require a valid JWT token in the Authorization header:
```
Authorization: Bearer your-jwt-token
```

To obtain a token (example flow):
1. Create a user via the `/api/v1/users` endpoint
2. Login via the `/api/v1/auth/login` endpoint to get a JWT token
3. Use this token for protected routes

## Environment Configuration

The application supports different environments through environment-specific .env files:
- `.env` - Default environment variables
- `.env.development` - Development environment overrides
- `.env.testing` - Testing environment overrides
- `.env.production` - Production environment overrides
- `.env.local` - Local overrides (not committed to git)

## Migrations

Migrations are automatically run when the application starts. You can also run them manually:
```bash
# Run migrations up
make migrate-up

# Roll back last migration
make migrate-down
```

## Docker Deployment

The included Dockerfile builds a production-ready container. Example deployment:
```bash
# Build the image
docker build -t fiber-gorm-api .

# Run the container
docker run -p 8080:8080 --env-file .env.production fiber-gorm-api
```

## CI/CD Pipeline

The project includes a GitHub Actions workflow for CI/CD in the `.github/workflows` directory. It:
1. Runs unit and integration tests
2. Performs code quality checks
3. Builds the application
4. Creates and pushes a Docker image
5. (Can be extended for deployment)

## Architecture

This application follows clean architecture principles:
1. **Handler Layer**: Handles HTTP requests and responses
2. **Service Layer**: Contains business logic
3. **Repository Layer**: Handles data access
4. **Model Layer**: Defines data structures

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Commit your changes: `git commit -am 'Add feature'`
4. Push to the branch: `git push origin feature-name`
5. Submit a pull request

## License

This project is licensed under the MIT License.