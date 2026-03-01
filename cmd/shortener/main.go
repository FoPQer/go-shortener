package main

import (
	"context"
	"net/http"

	"github.com/FoPQer/go-shortener/internal/config/db"
	"github.com/FoPQer/go-shortener/internal/config/flags"
	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/repository"
	"github.com/FoPQer/go-shortener/internal/routes"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	flags.ParseFlags()
	r := chi.NewRouter()
	routes.InitWebRoutes(r)
	logger.InitLogger()
	conn := db.InitPgsql()
	if conn != nil {
		defer conn.Close(context.Background())
	} else if service.GetFileStoragePath() != "" {
		repository.InitUrls(service.GetFileStoragePath())
	}

	if err := http.ListenAndServe(service.GetRunAddr(), r); err != nil {
		panic(err)
	}
}
