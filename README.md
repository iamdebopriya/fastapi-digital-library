# Digital Library - Dual Architecture Project

A comprehensive study of REST API architecture patterns through two parallel implementations: **FastAPI (Python)** and **Gin (Go)**. This project demonstrates best practices for building scalable, well-documented, and maintainable digital library management systems.

---

## Project Overview

This repository contains two functionally equivalent implementations of a digital book management API, allowing for direct architectural comparison between Python's asynchronous framework and Go's high-performance concurrency model.

**Common Features Across Both:**
- Complete CRUD operations for book management
- Request/response validation using domain models
- Automatic API documentation (Swagger/OpenAPI)
- Custom middleware for logging and performance tracking
- CORS support for cross-origin requests
- Background task simulation with request blocking
- In-memory data storage (no external database)

---

## File Structure & Components

### Backend Structure (FastAPI - Python)

```
backend/
├── main.py
│   ├── App Setup & Configuration
│   ├── Pydantic Models (validation layer)
│   ├── Global Middleware (logging, timing)
│   ├── CRUD Endpoints (5 core operations)
└── __pycache__/
```

**[backend/main.py](backend/main.py) - Core Implementation**
- **App Initialization**: FastAPI instance with title, description, and version
- **Data Validation**: Pydantic `Book` model with field validators
  - `year`: Range validation (1000–2026)
  - `isbn`: Length validation (10 or 13 characters)
- **Middleware**: Global logging middleware capturing User-Agent and processing time
- **Endpoints**: 5 REST endpoints for complete CRUD lifecycle
- **Storage**: In-memory list-based database (`library_db`)

### Backend Structure (Gin - Go)

```
digital-library-go/
├── cmd/
│   └── main.go (entry point, middleware config, server setup)
├── internal/
│   ├── delivery/http/
│   │   ├── handler.go (HTTP request handlers)
│   │   ├── router.go (route registration)
│   │   └── background_task_handler.go (task control)
│   └── usecase/
│       └── book_usecase.go (business logic, data storage)
└── docs/ (auto-generated Swagger documentation)
```

**[digital-library-go/cmd/main.go](digital-library-go/cmd/main.go) - Server Configuration**
- **Middleware Stack**: Three custom middleware functions
  - Task wait middleware (blocks requests if task running)
  - Timing + User-Agent middleware (performance tracking)
  - CORS middleware (cross-origin support)
- **Route Registration**: Delegates to `http.RegisterRoutes()`
- **Swagger Integration**: Automatic OpenAPI documentation

**[digital-library-go/internal/delivery/http/handler.go](digital-library-go/internal/delivery/http/handler.go) - Request Handlers**
- `GetBooks()` — Returns all books with 200 OK
- `GetBookByID()` — Retrieves single book or 404
- `CreateBook()` — Validates and persists new book; spawns background goroutine for email simulation
- `UpdateBook()` — Modifies existing book or returns 404
- `DeleteBook()` — Removes book from collection

**[digital-library-go/internal/usecase/book_usecase.go](digital-library-go/internal/usecase/book_usecase.go) - Business Logic**
- Manages in-memory book slice
- Implements CRUD operations with error handling
- Duplicate ID checking
- Dependency injection pattern via `BookHandler`

**[digital-library-go/internal/domain/book.go](digital-library-go/internal/domain/book.go) - Domain Model**
- Defines `Book` struct with JSON serialization tags
- `Validate()` method enforces constraints
- Identical validation rules to FastAPI version

**[digital-library-go/internal/delivery/http/background_task_handler.go](digital-library-go/internal/delivery/http/background_task_handler.go) - Task Management**
- Manages long-running background operations
- Uses mutex (`TaskMu`) to synchronize state
- Blocks concurrent task execution

**[digital-library-go/internal/delivery/http/router.go](digital-library-go/internal/delivery/http/router.go) - Route Registry**
- Centralized route definition
- Registers 6 endpoints: 5 CRUD + 1 background task

---

## Architecture Comparison

### FastAPI (Python)

```
Request
  ↓
Global Middleware (logging, timing)
  ↓
Route Handler (business logic embedded)
  ↓
Pydantic Validation
  ↓
In-Memory List Storage
```

**Characteristics:**
- **Monolithic Structure**: All logic in single `main.py`
- **Synchronous Processing**: Sequential request handling
- **Validation Layer**: Built-in Pydantic validators
- **Framework Features**: Auto-generated docs, easy middleware
- **Performance**: Good for I/O-bound tasks (async support)

### Gin (Go)

```
Request
  ↓
Middleware Stack
  ├─ Task Waiting
  ├─ Timing & Logging
  └─ CORS
  ↓
Router → Handler (HTTP layer)
  ↓
Usecase → Business Logic
  ↓
Domain → Data Validation
  ↓
In-Memory Slice Storage
```

**Characteristics:**
- **Layered Architecture**: Separation of concerns (delivery/usecase/domain)
- **Goroutine Support**: Concurrent request handling, background tasks
- **Explicit Validation**: Dedicated domain validation methods
- **Manual Documentation**: Swagger comments in code
- **Performance**: Exceptional for CPU-bound and concurrent operations
- **Type Safety**: Compile-time error detection

### Key Differences

| Aspect | FastAPI | Gin |
|--------|---------|-----|
| **Code Organization** | Single file | Layered (delivery/usecase/domain) |
| **Concurrency** | Async/await | Goroutines (lightweight) |
| **Request Handling** | Async-native | Sync (goroutine per request) |
| **Validation** | Pydantic decorators | Domain methods |
| **Documentation** | Auto-generated | Manual comments |
| **Middleware** | Function decorators | Explicit stack |
| **Startup Time** | Slower | Near-instantaneous |
| **Memory Usage** | Higher | Minimal |
| **Production Readiness** | High | High |

---

## API Endpoints

All endpoints are identical between implementations:

### Books CRUD

```bash
# Get all books
GET /books
Response: 200 {"data": [Book]}

# Get book by ID
GET /books/{id}
Response: 200 {"data": Book} or 404

# Create book
POST /books
Body: {"id": 1, "title": "...", "author": "...", "year": 2024, "isbn": "1234567890"}
Response: 201 or 400 (validation error)

# Update book
PUT /books/{id}
Body: {"id": 1, "title": "...", "author": "...", "year": 2024, "isbn": "..."}
Response: 200 or 404

# Delete book
DELETE /books/{id}
Response: 200 or 404
```

### Background Task

```bash
# Run blocking task
POST /tasks/process
Response: 200 {"message": "Task completed successfully"}
         or 409 {"error": "task already running"}
```

---

## Testing Examples

### FastAPI Testing

```python
# Start server
uvicorn backend.main:app --reload

# Test with curl
curl -X POST http://localhost:8000/books \
  -H "Content-Type: application/json" \
  -d '{"id": 1, "title": "1984", "author": "Orwell", "year": 1949, "isbn": "9780451524935"}'

# Test invalid year (should fail)
curl -X POST http://localhost:8000/books \
  -H "Content-Type: application/json" \
  -d '{"id": 2, "title": "Future Book", "author": "Author", "year": 2030, "isbn": "1234567890"}'
# Response: 422 Unprocessable Entity
```

### Gin Testing

```bash
# Start server
go run cmd/main.go

# Test create book
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"id": 1, "title": "Dune", "author": "Herbert", "year": 1965, "isbn": "0441172717"}'
# Response: 201 {"message": "book created"}

# Test invalid ISBN (should fail)
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"id": 2, "title": "Book", "author": "Author", "year": 2024, "isbn": "123"}'
# Response: 400 {"error": "isbn must be 10 or 13 characters"}

# Test background task
curl -X POST http://localhost:8080/tasks/process
# Response: 200 {"message": "Task completed successfully"}

# Test concurrent requests (task blocks others)
curl -X POST http://localhost:8080/tasks/process &
sleep 1
curl -X GET http://localhost:8080/books
# Second request will wait until task completes (8 seconds)
```

### Cross-Platform Testing

Both implementations return identical responses for the same requests:

```bash
# Valid book creation (both should succeed)
curl -X POST http://localhost:PORT/books \
  -H "Content-Type: application/json" \
  -d '{"id": 1, "title": "Foundation", "author": "Asimov", "year": 1951, "isbn": "0553382578"}'

# Duplicate ID check
curl -X POST http://localhost:PORT/books \
  -H "Content-Type: application/json" \
  -d '{"id": 1, "title": "Different Book", "author": "Different Author", "year": 2000, "isbn": "9876543210"}'
# Both: 400 Duplicate ID Error

# Retrieve all
curl -X GET http://localhost:PORT/books
# Both: 200 {"data": [...]}
```

---

## Architecture Diagram

### Request Flow (Unified)

```
┌─────────────────────────────────────────────────────────────────┐
│                       HTTP Request                              │
│              (GET /books, POST /books, etc.)                    │
└────────────────────────┬────────────────────────────────────────┘
                         │
         ┌───────────────┴───────────────┐
         ▼                               ▼
    ┌─────────────┐              ┌──────────────┐
    │   FastAPI   │              │     Gin      │
    │  (Python)   │              │     (Go)     │
    └─────────────┘              └──────────────┘
         │                               │
    ┌────▼──────────────┐          ┌────▼────────────────┐
    │   Middleware      │          │   Middleware Stack  │
    │ • Logging         │          │ • Task Wait         │
    │ • X-Process-Time  │          │ • Timing & Logging  │
    │ • CORS            │          │ • CORS              │
    └────▼──────────────┘          └────▼────────────────┘
         │                               │
    ┌────▼──────────────┐          ┌────▼────────────────┐
    │  FastAPI Route    │          │  Gin Handler        │
    │  (Decorator-based)│          │  (Method-based)     │
    └────▼──────────────┘          └────▼────────────────┘
         │                               │
    ┌────▼──────────────┐          ┌────▼────────────────┐
    │  Pydantic         │          │  Usecase            │
    │  Validation       │          │  Business Logic     │
    └────▼──────────────┘          └────▼────────────────┘
         │                               │
         └───────────────┬───────────────┘
                         ▼
         ┌───────────────────────────────┐
         │  Domain Validation            │
         │  • Year (1000–2026)           │
         │  • ISBN (10 or 13)            │
         │  • Title (non-empty)          │
         └───────────────┬───────────────┘
                         ▼
         ┌───────────────────────────────┐
         │  In-Memory Storage            │
         │  (List/Slice)                 │
         └───────────────┬───────────────┘
                         ▼
         ┌───────────────────────────────┐
         │  HTTP Response (200/400/404)  │
         └───────────────────────────────┘
```

### Concurrency Model

**FastAPI (Event Loop)**
```
Main Event Loop
    ├─ Request 1 (await async handler)
    ├─ Request 2 (await async handler)
    ├─ Request 3 (await async handler)
    └─ Context switches during I/O
```

**Gin (Goroutine Pool)**
```
Main Goroutine
    ├─ Goroutine 1 (Request handler)
    ├─ Goroutine 2 (Request handler)
    ├─ Goroutine 3 (Request handler)
    └─ OS scheduler manages concurrency
```

---

## Setup & Execution

### FastAPI Setup

```bash
# Navigate to project
cd fastapi-digital-library

# Create virtual environment
python -m venv .venv
.venv\Scripts\activate  # Windows

# Install dependencies
pip install -r requirements.txt

# Run server
uvicorn backend.main:app --reload --port 8000

# Access API docs
http://localhost:8000/docs
```

### Gin Setup

```bash
# Navigate to Go project
cd fastapi-digital-library/digital-library-go

# Download dependencies
go mod download
go mod tidy

# Run server
go run cmd/main.go

# Access API docs
http://localhost:8080/swagger/index.html
```

---

## Summary

This dual-implementation project demonstrates:

1. **Identical Functionality** — Both APIs provide the same features and behavior
2. **Different Paradigms** — Python's async simplicity vs. Go's goroutine efficiency
3. **Architecture Patterns** — Monolithic (FastAPI) vs. Layered (Gin)
4. **Developer Experience** — Auto-magic (Pydantic) vs. Explicit (Go domain validation)
5. **Performance Trade-offs** — Python's flexibility vs. Go's raw speed and resource efficiency

**Choose FastAPI for:** Rapid prototyping, data science integration, maximum developer productivity
**Choose Gin for:** High-concurrency scenarios, microservices, minimal resource consumption

---
Build with FastAPI(Python) & Gin(Go)
