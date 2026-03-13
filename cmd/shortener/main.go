package main

import (
	"context"
	"net/http"

	"github.com/FoPQer/go-shortener/internal/config/db"
	"github.com/FoPQer/go-shortener/internal/config/flags"
	"github.com/FoPQer/go-shortener/internal/handlers"
	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/repository/factory"
	"github.com/FoPQer/go-shortener/internal/routes"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	flags.ParseFlags()
	logger.InitLogger()
	pgxConf, err := db.InitPgsql()
	if err != nil {
		panic(err)
	}
	if pgxConf != nil {
		defer pgxConf.GetDBConn().Close()
	} 
		
    urlRepo, err := factory.
		NewRepositoryFactory(pgxConf.GetDBConn(), service.GetFileStoragePath()).
		CreateUrlsRepository(context.Background())
    if err != nil {
        panic(err)
    }
    
    urlService := service.NewURLService(urlRepo)
	jsonService := service.NewJSONService()

	handler := handlers.NewHandler(urlService, jsonService)
	dbHandler := handlers.NewDBHandler(pgxConf)

	r := chi.NewRouter()
	routes.InitWebRoutes(r, handler, dbHandler)

	if err := http.ListenAndServe(service.GetRunAddr(), r); err != nil {
		panic(err)
	}
}
