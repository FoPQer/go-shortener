package main

import (
	"net/http"
	"net/url"

	"github.com/FoPQer/go-shortener/internal/config/flags"
	"github.com/FoPQer/go-shortener/internal/handler"
	"github.com/FoPQer/go-shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	flags.ParseFlags()

	repository.InitUrls()
	r := chi.NewRouter()

	getPattern, err := url.JoinPath("/", flags.GetFlagBasePrefix())

	if err != nil {
		panic(err)
	}
	r.Get(getPattern+"/{id}", handler.GetURL)
	r.Post("/", handler.PostURL)

	if err := http.ListenAndServe(flags.GetFlagRunAddr(), r); err != nil {
		panic(err)
	}
}
