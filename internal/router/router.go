package router

import (
	"net/http"

	"desent-api-quest/internal/handler"
)

type AuthMiddleware func(http.Handler) http.Handler

func New(
	ping *handler.PingHandler,
	echo *handler.EchoHandler,
	auth *handler.AuthHandler,
	books *handler.BookHandler,
	requireBooksAuth AuthMiddleware,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", ping.Handle)
	mux.HandleFunc("POST /echo", echo.Handle)
	mux.HandleFunc("POST /auth/token", auth.IssueToken)
	mux.Handle("GET /books", requireBooksAuth(http.HandlerFunc(books.List)))
	mux.HandleFunc("POST /books", books.Create)
	mux.HandleFunc("GET /books/{id}", books.GetByID)
	mux.HandleFunc("PUT /books/{id}", books.Update)
	mux.HandleFunc("DELETE /books/{id}", books.Delete)

	return mux
}
