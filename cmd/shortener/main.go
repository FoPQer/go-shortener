package main

import (
	"context"
	"net/http"
	"os"

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
	repository.InitUrls(service.GetFileStoragePath())
	logger.InitLogger()
	err := db.InitPgsql()
	if err != nil {
		os.Exit(1)
		panic(err)
	}
	defer db.GetDBConn().Close(context.Background())

	if err := http.ListenAndServe(service.GetRunAddr(), r); err != nil {
		panic(err)
	}
}
