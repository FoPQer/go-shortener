package main

import (
	"net/http"

	"github.com/FoPQer/go-shortener/internal/handler"
	"github.com/FoPQer/go-shortener/internal/repository"
)

func main() {
	repository.InitUrls()
	mainPage := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetURL(w, r)
		case http.MethodPost:
			handler.PostURL(w, r)
		}
	}
	if err := http.ListenAndServe(":8080", http.HandlerFunc(mainPage)); err != nil {
		panic(err)
	}
}
