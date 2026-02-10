# Digital Library Go

A RESTful API for managing a digital book collection, built with Go and the Gin web framework. Provides comprehensive CRUD operations with built-in validation, interactive API documentation, and background task support.

## Overview

Digital Library Go is a lightweight, in-memory book management system designed to demonstrate clean architecture principles and API best practices in Go. It offers a fully functional REST API with OpenAPI/Swagger documentation, making it easy to explore and integrate.

## Features

- **Complete CRUD Operations** — Create, read, update, and delete books with full HTTP support
- **Data Validation** — Enforces constraints on book title, publication year (1000–2026), and ISBN format
- **Interactive Documentation** — Swagger UI available at `/swagger/` for real-time API exploration
- **CORS Support** — Enabled for cross-origin requests with proper response headers
- **Request Metrics** — Automatic `X-Process-Time` header on all responses for performance monitoring
- **Background Tasks** — Dedicated endpoint for simulating long-running operations

## Architecture

The application follows a layered architecture pattern with clear separation of concerns:

```
HTTP Requests
    ↓
HTTP Handlers (delivery/http)
    ↓
Business Logic (usecase)
    ↓
Domain Models (domain)
    ↓
In-Memory Storage
```

### Key Components

- `cmd/main.go` — Application entry point; configures middleware and routes
- `internal/delivery/http/handler.go` — HTTP request/response handlers for book operations
- `internal/usecase/book_usecase.go` — Core business logic and data storage
- `internal/domain/book.go` — `Book` data structure with validation logic

## Getting Started

### Prerequisites

- Go 1.16 or later
- `go` command-line tools

### Installation & Running

1. **Install dependencies:**
   ```bash
   go mod download
   go mod tidy
   ```

2. **Start the server:**
   ```bash
   go run cmd/main.go
   ```
   The API will be available at `http://localhost:8080`

3. **Access API documentation:**
   Open your browser to `http://localhost:8080/swagger/index.html#/`

## API Reference

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/books` | Retrieve all books |
| `GET` | `/books/:id` | Retrieve a specific book by ID |
| `POST` | `/books` | Create a new book (JSON body required) |
| `PUT` | `/books/:id` | Update an existing book (JSON body required) |
| `DELETE` | `/books/:id` | Delete a book by ID |
| `POST` | `/tasks/process` | Execute a background task simulation |

### Response Handling

- Successful operations return the appropriate HTTP 2xx status code with JSON data
- Validation errors return `400 Bad Request` with error details: `{"error": "..."}`
- Not found errors return `404 Not Found`

## Notes

- Data is stored in memory .
- No external database or persistent storage is used.
- Use the Swagger UI, `curl`, or your preferred HTTP client to test endpoints.

---

Built with Go and Gin
