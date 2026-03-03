package app

import (
	"net/http"

	"desent-api-quest/internal/handler"
	"desent-api-quest/internal/repository"
	"desent-api-quest/internal/router"
	"desent-api-quest/internal/token"
	"desent-api-quest/internal/usecase"
)

type App struct {
	addr    string
	handler http.Handler
}

func New() *App {
	bookRepo := repository.NewMemoryBookRepository()
	tokenSvc := token.NewService()

	bookUsecase := usecase.NewBookUsecase(bookRepo)
	authUsecase := usecase.NewAuthUsecase(tokenSvc)

	pingHandler := handler.NewPingHandler()
	echoHandler := handler.NewEchoHandler()
	authHandler := handler.NewAuthHandler(authUsecase)
	bookHandler := handler.NewBookHandler(bookUsecase)

	return &App{
		addr: ":8080",
		handler: router.New(
			pingHandler,
			echoHandler,
			authHandler,
			bookHandler,
		),
	}
}

func (a *App) Addr() string {
	return a.addr
}

func (a *App) Handler() http.Handler {
	return a.handler
}
