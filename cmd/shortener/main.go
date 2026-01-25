package main

import (
	"fmt"
	"net/http"

	"github.com/FoPQer/go-shortener/internal/config/flags"
	"github.com/FoPQer/go-shortener/internal/handler"
	"github.com/FoPQer/go-shortener/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	flags.ParseFlags()
	repository.InitUrls()
	r := chi.NewRouter()

	r.Route("/"+flags.FlagBasePrefix, func(r chi.Router) {
		r.Get("/{id}", handler.GetURL)
	})
	r.Post("/", handler.PostURL)
	fmt.Println(flags.FlagRunAddr + "/" + flags.FlagBasePrefix)
	if err := http.ListenAndServe(flags.FlagRunAddr, r); err != nil {
		panic(err)
	}
}
