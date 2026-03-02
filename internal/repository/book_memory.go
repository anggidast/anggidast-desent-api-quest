package repository

import (
	"strconv"
	"sync"

	"desent-api-quest/internal/domain"
)

type BookRepository interface {
	Create(input domain.BookInput) (domain.Book, error)
	List() ([]domain.Book, error)
	GetByID(id string) (domain.Book, error)
	Update(id string, input domain.BookInput) (domain.Book, error)
	Delete(id string) error
}

type MemoryBookRepository struct {
	mu      sync.RWMutex
	nextID  int
	books   map[string]domain.Book
	ordered []string
}

func NewMemoryBookRepository() *MemoryBookRepository {
	return &MemoryBookRepository{
		nextID: 1,
		books:  make(map[string]domain.Book),
	}
}

func (r *MemoryBookRepository) Create(input domain.BookInput) (domain.Book, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := strconv.Itoa(r.nextID)
	r.nextID++

	book := domain.Book{
		ID:     id,
		Title:  input.Title,
		Author: input.Author,
		Year:   input.Year,
	}

	r.books[id] = book
	r.ordered = append(r.ordered, id)

	return book, nil
}

func (r *MemoryBookRepository) List() ([]domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	books := make([]domain.Book, 0, len(r.ordered))
	for _, id := range r.ordered {
		book, ok := r.books[id]
		if ok {
			books = append(books, book)
		}
	}

	return books, nil
}

func (r *MemoryBookRepository) GetByID(id string) (domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, ok := r.books[id]
	if !ok {
		return domain.Book{}, &domain.FieldError{Code: domain.ErrNotFound, Message: "book not found"}
	}

	return book, nil
}

func (r *MemoryBookRepository) Update(id string, input domain.BookInput) (domain.Book, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.books[id]; !ok {
		return domain.Book{}, &domain.FieldError{Code: domain.ErrNotFound, Message: "book not found"}
	}

	book := domain.Book{
		ID:     id,
		Title:  input.Title,
		Author: input.Author,
		Year:   input.Year,
	}

	r.books[id] = book
	return book, nil
}

func (r *MemoryBookRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.books[id]; !ok {
		return &domain.FieldError{Code: domain.ErrNotFound, Message: "book not found"}
	}

	delete(r.books, id)
	return nil
}
