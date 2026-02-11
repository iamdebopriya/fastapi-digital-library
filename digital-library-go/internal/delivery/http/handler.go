package http

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/iamdebopriya/fastapi-digital-library/digital-library-go/internal/domain"
	"github.com/iamdebopriya/fastapi-digital-library/digital-library-go/internal/usecase"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	uc *usecase.BookUsecase
}

func NewBookHandler(uc *usecase.BookUsecase) *BookHandler {
	return &BookHandler{uc: uc}
}

// GetBooks godoc
// @Summary Get all books
// @Description Get list of all books
// @Tags Library
// @Produce json
// @Success 200 {array} domain.Book
// @Router /books [get]
func (h *BookHandler) GetBooks(c *gin.Context) {
	books := h.uc.GetBooks()
	c.JSON(http.StatusOK, gin.H{"data": books})
}

// GetBookByID godoc
// @Summary Get a book by ID
// @Description Get book details by ID
// @Tags Library
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} domain.Book
// @Failure 404 {object} map[string]string
// @Router /books/{id} [get]
func (h *BookHandler) GetBookByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	book, err := h.uc.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// CreateBook godoc
// @Summary Create a new book
// @Description Add a new book to the library
// @Tags Library
// @Accept json
// @Produce json
// @Param book body domain.Book true "Book data"
// @Success 201 {object} domain.Book
// @Failure 400 {object} map[string]string
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var book domain.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if err := book.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.CreateBook(book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	go func(b domain.Book) {
		time.Sleep(2 * time.Second)
		log.Println("Notification sent for new book:", b.Title)
	}(book)

	c.JSON(http.StatusCreated, gin.H{"message": "book created"})
}

// UpdateBook godoc
// @Summary Update a book
// @Description Update book details by ID
// @Tags Library
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body domain.Book true "Updated book data"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var book domain.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if err := book.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.uc.UpdateBook(id, book)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book updated"})
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Delete book by ID
// @Tags Library
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.uc.DeleteBook(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book deleted"})
}
