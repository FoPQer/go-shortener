package main

import (
	"net/http"

	"github.com/FoPQer/go-shortener/internal/config/flags"
	"github.com/FoPQer/go-shortener/internal/handler"
	"github.com/FoPQer/go-shortener/internal/repository"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	flags.ParseFlags()

	repository.InitUrls()
	r := chi.NewRouter()

	base := service.GetBasePrefix()

	r.Get(base+"{id}", handler.GetURL)
	r.Post("/", handler.PostURL)

	if err := http.ListenAndServe(service.GetRunAddr(), r); err != nil {
		panic(err)
	}
}
