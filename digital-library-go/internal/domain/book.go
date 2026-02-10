package domain

import "errors"

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	ISBN   string `json:"isbn"`
}

func (b *Book) Validate() error {
	if b.Title == "" {
		return errors.New("title must not be empty")
	}
	if b.Year < 1000 || b.Year > 2026 {
		return errors.New("year must be between 1000 and 2026")
	}
	if len(b.ISBN) != 10 && len(b.ISBN) != 13 {
		return errors.New("isbn must be 10 or 13 characters")
	}
	return nil
}
