package usecase

import (
	"errors"

	"github.com/iamdebopriya/fastapi-digital-library/digital-library-go/internal/domain"
)

type BookUsecase struct {
	books []domain.Book
}

func NewBookUsecase() *BookUsecase {
	return &BookUsecase{
		books: []domain.Book{},
	}
}

func (u *BookUsecase) GetBooks() []domain.Book {
	return u.books
}

func (u *BookUsecase) GetBookByID(id int) (domain.Book, error) {
	for _, b := range u.books {
		if b.ID == id {
			return b, nil
		}
	}
	return domain.Book{}, errors.New("book not found")
}

func (u *BookUsecase) isDuplicateID(id int) bool {
	for _, b := range u.books {
		if b.ID == id {
			return true
		}
	}
	return false
}

func (u *BookUsecase) CreateBook(book domain.Book) error {

	if u.isDuplicateID(book.ID) {
		return errors.New("book with this ID already exists")
	}

	u.books = append(u.books, book)
	return nil
}

func (u *BookUsecase) UpdateBook(id int, updated domain.Book) error {
	for i, b := range u.books {
		if b.ID == id {
			updated.ID = id
			u.books[i] = updated
			return nil
		}
	}
	return errors.New("book not found")
}

func (u *BookUsecase) DeleteBook(id int) error {
	for i, b := range u.books {
		if b.ID == id {
			u.books = append(u.books[:i], u.books[i+1:]...)
			return nil
		}
	}
	return errors.New("book not found")
}
