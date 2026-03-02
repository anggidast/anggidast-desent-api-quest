package handler

import (
	"net/http"
	"strconv"

	"desent-api-quest/internal/domain"
	"desent-api-quest/internal/httpjson"
	"desent-api-quest/internal/usecase"
)

type BookHandler struct {
	books *usecase.BookUsecase
}

func NewBookHandler(books *usecase.BookUsecase) *BookHandler {
	return &BookHandler{books: books}
}

func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.BookInput
	if err := httpjson.DecodeJSON(r, &input); err != nil {
		httpjson.WriteError(w, err)
		return
	}

	book, err := h.books.Create(input)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, book)
}

func (h *BookHandler) List(w http.ResponseWriter, r *http.Request) {
	page, err := parsePositiveIntOrDefault(r.URL.Query().Get("page"), 1, "page")
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}

	limit, err := parsePositiveIntOrDefault(r.URL.Query().Get("limit"), 10, "limit")
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}

	result, err := h.books.List(domain.ListBooksParams{
		Author: r.URL.Query().Get("author"),
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, result)
}

func (h *BookHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	book, err := h.books.GetByID(r.PathValue("id"))
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, book)
}

func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	var input domain.BookInput
	if err := httpjson.DecodeJSON(r, &input); err != nil {
		httpjson.WriteError(w, err)
		return
	}

	book, err := h.books.Update(r.PathValue("id"), input)
	if err != nil {
		httpjson.WriteError(w, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, book)
}

func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.books.Delete(r.PathValue("id")); err != nil {
		httpjson.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parsePositiveIntOrDefault(raw string, fallback int, name string) (int, error) {
	if raw == "" {
		return fallback, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, &domain.FieldError{Code: domain.ErrBadRequest, Message: name + " must be a valid integer"}
	}
	if value < 1 {
		return 0, &domain.FieldError{Code: domain.ErrBadRequest, Message: name + " must be greater than or equal to 1"}
	}

	return value, nil
}
