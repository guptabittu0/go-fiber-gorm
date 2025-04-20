# Go Fiber GORM API Boilerplate

A production-ready RESTful API boilerplate built with Go Fiber and GORM, designed to scale from small projects to large enterprise applications.

![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)
![Fiber Version](https://img.shields.io/badge/Fiber-v2.50%2B-blue)
![GORM Version](https://img.shields.io/badge/GORM-v1.25%2B-blue)
![License](https://img.shields.io/badge/License-MIT-green)

## 🌟 Key Features

- **[Go Fiber](https://github.com/gofiber/fiber)**: Ultra-fast HTTP framework built on top of Fasthttp
- **[GORM](https://gorm.io)**: Feature-rich ORM with PostgreSQL integration
- **Clean Architecture**: Well-structured code with separation of concerns
- **API Versioning**: Support for multiple API versions
- **Authentication**: Complete JWT-based auth system with refresh tokens
- **Authorization**: Role-based access control for fine-grained permissions
- **Middleware**: Logging, error handling, rate limiting, JWT validation
- **Configuration**: Environment-based configs with `.env` file support
- **Database Migrations**: Automatic and manual migration support
- **Redis Integration**: Caching and session management
- **Validation**: Request validation using validator tags
- **Structured Error Handling**: Consistent error responses
- **Health Checks**: System monitoring endpoints
- **Background Processing**: Worker pool for async tasks
- **Docker Support**: Containerized development and deployment
- **CI/CD Pipeline**: GitHub Actions workflow
- **Hot Reload**: Fast development with Air
- **Comprehensive Testing**: Framework for unit and integration tests

## 📁 Project Structure

```
├── cmd/                          # Application entry points
│   ├── main.go                   # Main application
│   └── migrate/                  # Database migration tool
├── config/                       # Configuration
│   ├── config.go                 # Configuration structs
│   └── env_loader.go             # Environment loader
├── core/                         # Core framework components
│   ├── cache/                    # Redis integration
│   ├── database/                 # Database connection/transaction
│   ├── errors/                   # Error handling
│   ├── logger/                   # Logging utilities
│   ├── middleware/               # Global middleware
│   └── worker/                   # Background worker pool
├── migrations/                   # Migration definitions
├── modules/                      # Feature modules
│   ├── auth/                     # Authentication/authorization
│   ├── health/                   # Health check endpoints
│   └── user/                     # User management
├── routes/                       # Route registration
├── test/                         # Testing utilities
├── docker-compose.yml            # Docker services
├── Dockerfile                    # Container definition
├── go.mod                        # Dependencies
└── Makefile                      # Build commands
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- Redis (optional but recommended)
- Docker and Docker Compose (optional)

### Local Development

1. Clone the repository
   ```bash
   git clone https://github.com/yourusername/go-fiber-gorm.git
   cd go-fiber-gorm
   ```

2. Copy and configure environment variables
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. Install dependencies
   ```bash
   go mod tidy
   ```

4. Run database migrations
   ```bash
   go run cmd/migrate/migrate.go
   ```

5. Run the application
   ```bash
   go run cmd/main.go
   ```

6. For hot reloading during development
   ```bash
   make dev
   # or directly: air
   ```

### Docker Development

```bash
# Start all services (app, postgres, redis)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

## 📋 Available Make Commands

```bash
# Build the application
make build

# Run the application
make run

# Run with hot reload (using air)
make dev

# Run tests
make test

# Run with test coverage
make coverage

# Run database migrations
make migrate-up

# Rollback last migration
make migrate-down

# Build Docker image
make docker-build

# Deploy with Docker Compose
make docker-compose-up

# Generate a new module
make gen-module
# You'll be prompted for the module name
```

## 🔐 Authentication

The boilerplate includes a complete JWT authentication system:

1. Register a new user:
   ```http
   POST /api/v1/auth/register
   {
     "name": "John Doe",
     "email": "john@example.com",
     "password": "securepassword"
   }
   ```

2. Login to get access and refresh tokens:
   ```http
   POST /api/v1/auth/login
   {
     "email": "john@example.com",
     "password": "securepassword"
   }
   ```

3. Use the JWT token in subsequent requests:
   ```http
   GET /api/v1/users
   Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
   ```

4. Refresh your token when it expires:
   ```http
   POST /api/v1/auth/refresh-token
   {
     "refresh_token": "your-refresh-token"
   }
   ```

## 🏗 API Routes

### Auth Module
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh-token` - Refresh access token
- `POST /api/v1/auth/logout` - Logout (invalidate current session)
- `POST /api/v1/auth/logout-all` - Logout from all devices
- `POST /api/v1/auth/change-password` - Change user password

### User Module
- `POST /api/v1/users` - Create a user (admin only)
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user (admin only)

### Health Module
- `GET /api/v1/health` - Basic health check
- `GET /api/v1/health/details` - Detailed health check with component status

## ⚙️ Configuration

The application uses environment-specific configuration files:

- `.env` - Default environment variables
- `.env.development` - Development-specific overrides
- `.env.testing` - Testing overrides
- `.env.production` - Production overrides
- `.env.local` - Local overrides (not committed to git)

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Port for the HTTP server | `8080` |
| `SERVER_ENV` | Environment (development/production) | `development` |
| `SERVER_TIMEOUT` | Request timeout in seconds | `10` |
| `SERVER_READ_TIMEOUT` | Read timeout in seconds | `15` |
| `SERVER_WRITE_TIMEOUT` | Write timeout in seconds | `15` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database username | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `fiber_gorm` |
| `DB_SSL_MODE` | Database SSL mode | `disable` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `REDIS_PASSWORD` | Redis password | - |
| `JWT_SECRET` | Secret key for JWT | `your-secret-key` |
| `JWT_EXPIRY` | JWT expiration time | `15m` |
| `REFRESH_TOKEN_EXPIRY` | Refresh token expiration | `168h` |

## 🧪 Testing

The project includes utilities for both unit and integration tests:

```bash
# Run all tests
make test

# Run tests with coverage report
make coverage
```

For integration tests, the `TestRequest` utility simplifies API testing:

```go
func TestUserCreation(t *testing.T) {
    app := test.SetupTestApp()
    
    resp := test.MakeTestRequest(t, app, test.TestRequest{
        Method: "POST",
        URL:    "/api/v1/auth/register",
        Body: map[string]string{
            "name":     "Test User",
            "email":    "test@example.com",
            "password": "password123",
        },
    })
    
    assert.Equal(t, 201, resp.StatusCode)
}
```

## 🧱 Architecture

This boilerplate follows clean architecture principles:

1. **Controllers/Handlers**: Handle HTTP requests/responses
2. **Services**: Implement business logic
3. **Repositories**: Handle data access and storage
4. **Models**: Define data structures
5. **DTOs**: Handle data transfer objects

Each module (like `auth`, `user`) contains its own implementation of these components, making the codebase modular and maintainable.

## 🔧 Performance Optimizations

- Connection pooling for database and Redis
- Request rate limiting
- Efficient JSON serialization/deserialization
- Middleware execution optimization
- Database query optimization with GORM
- Redis-based caching for frequently accessed data
- Worker pool for background processing
- Middleware short-circuiting

## 🔒 Security Features

- JWT token-based authentication
- Role-based authorization
- Password hashing with bcrypt
- Request validation to prevent injection attacks
- Rate limiting to prevent brute force attacks
- CORS protection
- XSS protection headers
- SQL injection protection via GORM

## 🚢 Deployment

The project includes Docker and Docker Compose configurations for easy deployment:

```bash
# Build the Docker image
docker build -t go-fiber-gorm-api .

# Run the container
docker run -p 8080:8080 --env-file .env.production go-fiber-gorm-api
```

## 📚 Credits

- [Fiber](https://github.com/gofiber/fiber) - Express-inspired web framework
- [GORM](https://gorm.io) - ORM library for Go
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT implementation
- [Testify](https://github.com/stretchr/testify) - Testing toolkit
- [Air](https://github.com/cosmtrek/air) - Live reload for Go apps
- [Validator](https://github.com/go-playground/validator) - Request validation
- [Go-Redis](https://github.com/go-redis/redis) - Redis client for Go

## 📄 License

This project is licensed under the MIT License.

## 👥 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🙏 Acknowledgments

- Special thanks to the Go, Fiber, and GORM communities
- Inspired by best practices in Go API development
- Built with ❤️ for high-performance Go applications