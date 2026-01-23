# FastAPI Digital Library

A RESTful Digital Library API built using FastAPI. This project demonstrates proper Git workflow, strict data validation, middleware logging, error handling, and self-documenting APIs as part of TAG COE Training 2026 – Week 3 Assignment.

The application allows users to Create, Read, Update, and Delete (CRUD) books stored in an in-memory database.

---

## Features

- FastAPI-based RESTful API
- Pydantic v2 data modeling and validation
- Full CRUD operations
- Custom error handling (400 / 404 / 422)
- Global middleware for request logging
- Request processing time tracking
- Swagger / OpenAPI self-documentation

---

## Project Structure

```
fastapi-digital-library/
│
├── backend/
│   └── main.py              # FastAPI application entry point
│
├── requirements.txt         # Project dependencies
├── README.md                # Project documentation
└── .gitignore               # Python ignore rules
```

---

## Tech Stack

- Python 3.9+
- FastAPI
- Uvicorn
- Pydantic v2

---

## Setup Instructions

### Clone Repository

```bash
git clone https://github.com/iamdebopriya/fastapi-digital-library.git
cd fastapi-digital-library
```

### 1. Create Virtual Environment

```bash
python -m venv .venv
.venv\Scripts\activate      # Windows
source .venv/bin/activate   # macOS/Linux
```

### 2. Install Dependencies

```bash
pip install -r requirements.txt
```

### 3. Run the Application

```bash
uvicorn backend.main:app --reload
```

Server runs at: `http://127.0.0.1:8000`

---

## API Documentation (Swagger)

FastAPI provides automatic API documentation.

- **Swagger UI**: `http://127.0.0.1:8000/docs`
- **OpenAPI JSON**: `http://127.0.0.1:8000/openapi.json`

---

## Book Data Model

| Field  | Type | Constraints              |
|--------|------|--------------------------|
| id     | int  | Unique identifier        |
| title  | str  | Must not be empty        |
| author | str  | Required                 |
| year   | int  | 1000 – 2026              |
| isbn   | str  | Length must be 10 or 13  |

---

## API Endpoints

### 1. Create Book

**POST** `/books`

**Request Body:**

```json
{
  "id": 1,
  "title": "Clean Code",
  "author": "Robert C. Martin",
  "year": 2008,
  "isbn": "0132350882"
}
```

**Success Response (200):**

```json
{
  "id": 1,
  "title": "Clean Code",
  "author": "Robert C. Martin",
  "year": 2008,
  "isbn": "0132350882"
}
```

**Duplicate ID Error (400):**

```json
{
  "detail": "Duplicate ID Error: Book with this ID already exists"
}
```

---

### 2. Get All Books

**GET** `/books`

**Response (200):**

```json
[
  {
    "id": 1,
    "title": "Clean Code",
    "author": "Robert C. Martin",
    "year": 2008,
    "isbn": "0132350882"
  }
]
```

---

### 3. Get Book by ID

**GET** `/books/{book_id}`

**Example:** `/books/1`

**Response (200):**

```json
{
  "id": 1,
  "title": "Clean Code",
  "author": "Robert C. Martin",
  "year": 2008,
  "isbn": "0132350882"
}
```

**Not Found (404):**

```json
{
  "detail": "Resource Not Found: Book with this ID does not exist"
}
```

---

### 4. Update Book

**PUT** `/books/{book_id}`

**Request Body:**

```json
{
  "id": 1,
  "title": "Clean Code (Updated)",
  "author": "Robert C. Martin",
  "year": 2009,
  "isbn": "0132350882"
}
```

**Response (200):**

```json
{
  "id": 1,
  "title": "Clean Code (Updated)",
  "author": "Robert C. Martin",
  "year": 2009,
  "isbn": "0132350882"
}
```

**Not Found (404):**

```json
{
  "detail": "Resource Not Found: Cannot update non-existent book"
}
```

---

### 5. Delete Book

**DELETE** `/books/{book_id}`

**Example:** `/books/1`

**Response (200):**

```json
{
  "message": "Book removed successfully"
}
```

**Not Found (404):**

```json
{
  "detail": "Resource Not Found: Cannot delete non-existent book"
}
```

---

## Validation Errors (422)

Returned automatically by Pydantic.

**Example:**

```json
{
  "detail": [
    {
      "type": "value_error",
      "loc": [
        "body",
        "year"
      ],
      "msg": "Value error, Year must be between 1000 and 2026",
      "input": 3009,
      "ctx": {
        "error": {}
      }
    }
  ]
}

{
  "detail": [
    {
      "type": "value_error",
      "loc": [
        "body",
        "isbn"
      ],
      "msg": "Value error, ISBN must be 10 or 13 characters long",
      "input": "12",
      "ctx": {
        "error": {}
      }
    }
  ]
}
```

---

## Middleware Logging

A global middleware:

- Captures User-Agent
- Logs incoming requests
- Measures request processing time
- Adds `X-Process-Time` header

**Example console output:**

```
[LOG] Request received from: Mozilla/5.0
```

---

## Git Workflow

- Feature branch: `digital-library-service`
- Changes pushed to GitHub
- Pull Request created
- Code reviewed and merged into `main`

---

## License

This project is licensed under the MIT License.

---

## Author

**Debopriya Lahiri**

TAG COE Training – 2026 Week 3 Assignment