package main

import (
	"net/http"
	"strings"

	"github.com/FoPQer/go-shortener/internal/config/flags"
	"github.com/FoPQer/go-shortener/internal/handler"
	"github.com/FoPQer/go-shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	flags.ParseFlags()

	repository.InitUrls()
	r := chi.NewRouter()

	base := flags.GetFlagBasePrefix()
	if !strings.HasPrefix(base, "/") {
		base = "/" + base
	}
	if !strings.HasSuffix(base, "/") {
		base = base + "/"
	}

	r.Get(base+"{id}", handler.GetURL)
	r.Post("/", handler.PostURL)

	if err := http.ListenAndServe(flags.GetFlagRunAddr(), r); err != nil {
		panic(err)
	}
}
