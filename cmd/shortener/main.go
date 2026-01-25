package main

import (
	"net/http"

	"github.com/FoPQer/go-shortener/internal/handler"
	"github.com/FoPQer/go-shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	repository.InitUrls()
	r := chi.NewRouter()

	r.Post("/", handler.PostURL)
	r.Get("/{id}", handler.GetURL)
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
