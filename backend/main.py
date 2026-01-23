from fastapi import FastAPI, HTTPException, Request
from pydantic import BaseModel, Field, field_validator
from typing import List
import time


# App Initialization


app = FastAPI(
    title="Digital Library API",
    description="Build a RESTful API to manage a digital library. "
                "The system ensures data integrity through strict validation, "
                "provides clear documentation, and tracks incoming traffic metadata.",
    version="1.0.0"
)


# Middleware Logging (Global)


@app.middleware("http")
async def global_logging_middleware(request: Request, call_next):
    start_time = time.time()

    # Capture User-Agent
    user_agent = request.headers.get("user-agent", "Unknown")

    # Log to console
    print(f"[LOG] Request received from: {user_agent}")

    # Process request
    response = await call_next(request)

    # Calculate processing time
    process_time = time.time() - start_time

    # Add custom response header
    response.headers["X-Process-Time"] = str(process_time)

    return response


# Data Model (Pydantic)


class Book(BaseModel):
    """
    Book model used for request and response bodies.
    Acts as the blueprint for the Digital Library API.
    """
    id: int
    title: str = Field(..., min_length=1, description="Title must not be empty")
    author: str
    year: int
    isbn: str

    # Year validation (1000–2026)
    @field_validator("year")
    @classmethod
    def validate_year(cls, value: int):
        if value < 1000 or value > 2026:
            raise ValueError("Year must be between 1000 and 2026")
        return value

    # ISBN validation (10 or 13 chars)
    @field_validator("isbn")
    @classmethod
    def validate_isbn(cls, value: str):
        if len(value) not in (10, 13):
            raise ValueError("ISBN must be 10 or 13 characters long")
        return value


# In-Memory Database


library_db: List[Book] = []


# CRUD Operations (Grouped under 'Library' tag)



# Create

@app.post(
    "/books",
    tags=["Library"],
    summary="Add a new book",
    description="Add a new book to the in-memory digital library database"
)
def create_book(book: Book):
    # Duplicate ID check
    for existing_book in library_db:
        if existing_book.id == book.id:
            raise HTTPException(
                status_code=400,
                detail="Duplicate ID Error: Book with this ID already exists"
            )
    library_db.append(book)
    return book


# Read All

@app.get(
    "/books",
    tags=["Library"],
    summary="Retrieve all books",
    description="Retrieve all books stored in the digital library"
)
def read_all_books():
    return library_db


# Read One

@app.get(
    "/books/{book_id}",
    tags=["Library"],
    summary="Retrieve book by ID",
    description="Retrieve a specific book from the library using its ID"
)
def read_book(book_id: int):
    for book in library_db:
        if book.id == book_id:
            return book

    # Resource Not Found
    raise HTTPException(
        status_code=404,
        detail="Resource Not Found: Book with this ID does not exist"
    )


# Update

@app.put(
    "/books/{book_id}",
    tags=["Library"],
    summary="Update a book",
    description="Modify existing book details using the book ID"
)
def update_book(book_id: int, updated_book: Book):
    for index, book in enumerate(library_db):
        if book.id == book_id:
            library_db[index] = updated_book
            return updated_book

    # Resource Not Found
    raise HTTPException(
        status_code=404,
        detail="Resource Not Found: Cannot update non-existent book"
    )


# Delete

@app.delete(
    "/books/{book_id}",
    tags=["Library"],
    summary="Delete a book",
    description="Remove a book from the digital library system"
)
def delete_book(book_id: int):
    for book in library_db:
        if book.id == book_id:
            library_db.remove(book)
            return {"message": "Book removed successfully"}

    # Resource Not Found
    raise HTTPException(
        status_code=404,
        detail="Resource Not Found: Cannot delete non-existent book"
    )
