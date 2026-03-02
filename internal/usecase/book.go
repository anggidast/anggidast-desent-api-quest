package usecase

import (
	"strings"

	"desent-api-quest/internal/domain"
	"desent-api-quest/internal/repository"
)

type BookUsecase struct {
	repo repository.BookRepository
}

func NewBookUsecase(repo repository.BookRepository) *BookUsecase {
	return &BookUsecase{repo: repo}
}

func (u *BookUsecase) Create(input domain.BookInput) (domain.Book, error) {
	normalized, err := normalizeAndValidateBookInput(input)
	if err != nil {
		return domain.Book{}, err
	}

	return u.repo.Create(normalized)
}

func (u *BookUsecase) List(params domain.ListBooksParams) (domain.BookListResult, error) {
	page := params.Page
	if page == 0 {
		page = 1
	}
	limit := params.Limit
	if limit == 0 {
		limit = 10
	}
	if page < 1 {
		return domain.BookListResult{}, &domain.FieldError{Code: domain.ErrBadRequest, Message: "page must be greater than or equal to 1"}
	}
	if limit < 1 {
		return domain.BookListResult{}, &domain.FieldError{Code: domain.ErrBadRequest, Message: "limit must be greater than or equal to 1"}
	}

	books, err := u.repo.List()
	if err != nil {
		return domain.BookListResult{}, err
	}

	authorFilter := strings.TrimSpace(strings.ToLower(params.Author))
	filtered := make([]domain.Book, 0, len(books))
	for _, book := range books {
		if authorFilter != "" && strings.ToLower(strings.TrimSpace(book.Author)) != authorFilter {
			continue
		}
		filtered = append(filtered, book)
	}

	totalItems := len(filtered)
	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + limit - 1) / limit
	}

	start := (page - 1) * limit
	end := start + limit
	items := []domain.Book{}
	if start < totalItems {
		if end > totalItems {
			end = totalItems
		}
		items = filtered[start:end]
	}

	return domain.BookListResult{
		Items:      items,
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (u *BookUsecase) GetByID(id string) (domain.Book, error) {
	if strings.TrimSpace(id) == "" {
		return domain.Book{}, &domain.FieldError{Code: domain.ErrBadRequest, Message: "book id is required"}
	}
	return u.repo.GetByID(id)
}

func (u *BookUsecase) Update(id string, input domain.BookInput) (domain.Book, error) {
	if strings.TrimSpace(id) == "" {
		return domain.Book{}, &domain.FieldError{Code: domain.ErrBadRequest, Message: "book id is required"}
	}

	normalized, err := normalizeAndValidateBookInput(input)
	if err != nil {
		return domain.Book{}, err
	}

	return u.repo.Update(id, normalized)
}

func (u *BookUsecase) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return &domain.FieldError{Code: domain.ErrBadRequest, Message: "book id is required"}
	}
	return u.repo.Delete(id)
}

func normalizeAndValidateBookInput(input domain.BookInput) (domain.BookInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Author = strings.TrimSpace(input.Author)

	switch {
	case input.Title == "":
		return domain.BookInput{}, &domain.FieldError{Code: domain.ErrValidation, Message: "title is required"}
	case input.Author == "":
		return domain.BookInput{}, &domain.FieldError{Code: domain.ErrValidation, Message: "author is required"}
	case input.Year <= 0:
		return domain.BookInput{}, &domain.FieldError{Code: domain.ErrValidation, Message: "year must be greater than 0"}
	default:
		return input, nil
	}
}
