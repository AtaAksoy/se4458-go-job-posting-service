# Job Posting Service API

A production-ready RESTful API service for job posting management with Redis caching, built with Go, GORM, and Gin framework.

## 🚀 Features

- **CRUD Operations**: Complete Create, Read, Update, Delete operations for job postings
- **Search & Pagination**: Advanced search functionality with pagination support
- **Redis Caching**: High-performance caching layer for improved response times
- **API Versioning**: Clean API versioning structure (`/api/v1/`)
- **Swagger Documentation**: Interactive API documentation
- **Clean Architecture**: Well-structured, maintainable codebase
- **MySQL Support**: Production-ready database integration
- **Environment Configuration**: Flexible configuration management

## 🏗️ Architecture & Design

### Clean Architecture Approach
The project follows clean architecture principles with clear separation of concerns:

```
se4458-go-job-posting-service/
├── cmd/
│   └── main.go                 # Application entry point
├── config/
│   └── config.go               # Configuration management
├── internal/
│   ├── v1/                     # API version 1
│   │   ├── jobs/               # Job domain
│   │   │   ├── handler.go      # HTTP handlers
│   │   │   ├── repository.go   # Data access layer
│   │   │   ├── model.go        # Domain models
│   │   │   ├── dto.go          # Data Transfer Objects
│   │   │   └── cache.go        # Redis cache layer
│   │   └── db/
│   │       ├── db.go           # Database connection
│   │       └── redis.go        # Redis client
│   └── router.go               # Route definitions
├── docs/                       # Swagger documentation
├── .env                        # Environment variables
├── .gitignore
├── go.mod
└── README.md
```

### Design Patterns Used

1. **Repository Pattern**: Abstract data access layer
2. **DTO Pattern**: Separate request/response structures
3. **Cache-First Strategy**: Redis caching with fallback to database
4. **Dependency Injection**: Clean dependency management
5. **Middleware Pattern**: Request/response processing

## 📊 Data Models

### Job Entity
```go
type Job struct {
    ID          uint   `gorm:"primaryKey" json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Company     string `json:"company"`
    City        string `json:"city"`
    State       string `json:"state"`
    Status      bool   `json:"status"`
    CreatedAt   int64  `json:"created_at"`
}
```

### Entity Relationship Diagram (ERD)
```
┌─────────────────┐
│      Jobs       │
├─────────────────┤
│ id (PK)         │
│ title           │
│ description     │
│ company         │
│ city            │
│ state           │
│ status          │
│ created_at      │
└─────────────────┘
```

### DTOs (Data Transfer Objects)

#### CreateJobRequest
```go
type CreateJobRequest struct {
    Title       string `json:"title" binding:"required"`
    Description string `json:"description" binding:"required"`
    Company     string `json:"company" binding:"required"`
    City        string `json:"city" binding:"required"`
    State       string `json:"state" binding:"required"`
}
```

#### UpdateJobRequest (Partial Update)
```go
type UpdateJobRequest struct {
    Title       *string `json:"title"`
    Description *string `json:"description"`
    Company     *string `json:"company"`
    City        *string `json:"city"`
    State       *string `json:"state"`
    Status      *bool   `json:"status"`
}
```

#### JobResponse
```go
type JobResponse struct {
    ID          uint   `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Company     string `json:"company"`
    City        string `json:"city"`
    State       string `json:"state"`
    CreatedAt   int64  `json:"created_at"`
    Status      bool   `json:"status"`
}
```

## 🔧 API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

| Method | Endpoint | Description | Cache Strategy |
|--------|----------|-------------|----------------|
| POST | `/jobs` | Create a new job | Cache job, invalidate lists |
| GET | `/jobs` | List jobs with pagination | Cache lists (15min TTL) |
| GET | `/jobs/:id` | Get job by ID | Cache individual jobs (30min TTL) |
| PUT | `/jobs/:id` | Update job (partial) | Invalidate related caches |
| DELETE | `/jobs/:id` | Delete job | Invalidate all related caches |
| GET | `/jobs/search` | Search jobs | Cache search results (10min TTL) |

### Query Parameters
- `page`: Page number (default: 1)
- `limit`: Page size (default: 10)
- `q`: Search query (for search endpoint)

### Example Usage

#### Create Job
```bash
curl -X POST http://localhost:8080/api/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Senior Go Developer",
    "description": "We are looking for an experienced Go developer",
    "company": "TechCorp",
    "city": "Istanbul",
    "state": "TR"
  }'
```

#### List Jobs
```bash
curl "http://localhost:8080/api/v1/jobs?page=1&limit=10"
```

#### Search Jobs
```bash
curl "http://localhost:8080/api/v1/jobs/search?q=developer&page=1&limit=5"
```

## 🚀 Getting Started

### Prerequisites
- Go 1.19+
- MySQL 8.0+
- Redis 6.0+
- Git

### Option 1: Local Development

#### Installation

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/se4458-go-job-posting-service.git
cd se4458-go-job-posting-service
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**
Create a `.env` file in the root directory:
```env
DB_DSN=username:password@tcp(localhost:3306)/jobsdb?parseTime=true
PORT=8080
REDIS_ADDR=localhost:6379
REDIS_DB=0
REDIS_PASSWORD=
```

4. **Create MySQL database**
```sql
CREATE DATABASE jobsdb;
```

5. **Run the application**
```bash
go run cmd/main.go
```

6. **Access Swagger documentation**
```
http://localhost:8080/swagger/index.html
```

### Option 2: Docker Deployment

#### Quick Start with Docker Compose

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/se4458-go-job-posting-service.git
cd se4458-go-job-posting-service
```

2. **Start all services**
```bash
docker-compose up -d
```

3. **Access the application**
- API: http://localhost:8080/api/v1
- Swagger UI: http://localhost:8080/swagger/index.html
- Redis Commander: http://localhost:8081

#### Manual Docker Build

1. **Build the image**
```bash
docker build -t job-posting-service .
```

2. **Run the container**
```bash
docker run -p 8080:8080 \
  -e DB_DSN="root:password@tcp(host.docker.internal:3306)/jobsdb?parseTime=true" \
  -e REDIS_ADDR="host.docker.internal:6379" \
  job-posting-service
```

#### Docker Services

The `docker-compose.yml` includes:
- **app**: Go application (port 8080)
- **mysql**: MySQL 8.0 database (port 3306)
- **redis**: Redis 7 cache (port 6379)
- **redis-commander**: Redis management UI (port 8081)

#### Sample Data

The application comes with sample job data that will be automatically inserted when using Docker Compose.

## 🔄 Cache Strategy

### Redis Cache Implementation
- **Cache-First Approach**: Check cache before database
- **Automatic Invalidation**: Cache invalidation on data changes
- **TTL Management**: Different TTL for different data types
- **Pattern-Based Invalidation**: Bulk cache clearing for related data

### Cache Keys
- Individual jobs: `job:{id}`
- Job lists: `jobs:list:{page}:{limit}`
- Search results: `jobs:search:{query}:{page}:{limit}`

### Cache TTL
- Individual jobs: 30 minutes
- Job lists: 15 minutes
- Search results: 10 minutes

## 🛠️ Technologies Used

- **Go**: Programming language
- **Gin**: HTTP web framework
- **GORM**: ORM library
- **MySQL**: Primary database
- **Redis**: Caching layer
- **Swagger**: API documentation
- **godotenv**: Environment configuration

## 📝 Assumptions & Design Decisions

### Assumptions
1. **Job Status**: Jobs have a boolean status field for active/inactive
2. **Location Structure**: Location is split into city and state for better search
3. **Pagination**: Default page size of 10 items
4. **Cache Strategy**: Cache-first with database fallback
5. **API Versioning**: Current version is v1, future versions will be v2, v3, etc.

### Design Decisions
1. **Clean Architecture**: Separation of concerns for maintainability
2. **Repository Pattern**: Abstract data access for testability
3. **DTO Pattern**: Separate request/response structures for API stability
4. **Context Usage**: Proper context propagation for timeouts and cancellation
5. **Error Handling**: Consistent error responses across all endpoints
6. **Partial Updates**: Support for updating only specific fields

## 🧪 Testing

### Manual Testing
Use the Swagger UI at `http://localhost:8080/swagger/index.html` to test all endpoints.

### Example Test Scenarios
1. Create a job and verify it appears in the list
2. Search for jobs using different query terms
3. Update a job and verify changes are reflected
4. Delete a job and verify it's removed from all endpoints
5. Test pagination with different page sizes

## 📈 Performance Considerations

### Caching Benefits
- **Response Time**: Cached responses are 10-100x faster
- **Database Load**: Reduced database queries
- **Scalability**: Better handling of concurrent requests

### Optimization Strategies
- **Connection Pooling**: Database and Redis connection pooling
- **Indexing**: Proper database indexes on search fields
- **Lazy Loading**: Cache population on-demand
- **TTL Management**: Appropriate cache expiration times

## 🔒 Security Considerations

### Input Validation
- Required field validation using Gin binding
- SQL injection prevention through GORM
- Input sanitization for search queries

### Environment Security
- Sensitive data in environment variables
- Database credentials not hardcoded
- Redis password support

## 🚀 Deployment

### Docker Support
```dockerfile
FROM golang:1.19-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/main.go
EXPOSE 8080
CMD ["./main"]
```

### Environment Variables for Production
```env
DB_DSN=production_user:secure_password@tcp(db.example.com:3306)/jobsdb?parseTime=true
REDIS_ADDR=redis.example.com:6379
REDIS_PASSWORD=secure_redis_password
PORT=8080
```

## 📚 API Documentation

Interactive API documentation is available at:
```
http://localhost:8080/swagger/index.html
```

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Go community for excellent documentation
- Gin framework for the robust HTTP framework
- GORM for the powerful ORM library
- Redis for the fast caching solution 