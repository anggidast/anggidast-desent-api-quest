package main

import (
	"log"
	"net/http"

	"desent-api-quest/internal/app"
)

func main() {
	application := app.New()

	log.Printf("starting server on %s", application.Addr())
	if err := http.ListenAndServe(application.Addr(), application.Handler()); err != nil {
		log.Fatal(err)
	}
}
