package main

import (
	"context"
	"net/http"

	"github.com/FoPQer/go-shortener/internal/config/db"
	"github.com/FoPQer/go-shortener/internal/config/flags"
	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/repository/factory"
	"github.com/FoPQer/go-shortener/internal/routes"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	flags.ParseFlags()
	logger.InitLogger()
	conn := db.InitPgsql()
	if conn != nil {
		defer conn.Close(context.Background())
	} 
		
    urlRepo, err := factory.
		NewRepositoryFactory(conn, service.GetFileStoragePath()).
		CreateUrlsRepository(context.Background())
    if err != nil {
        panic(err)
    }
    
    urlService := service.NewURLService(urlRepo)
	jsonService := service.NewJSONService()

	r := chi.NewRouter()
	routes.InitWebRoutes(r, urlService, jsonService)

	if err := http.ListenAndServe(service.GetRunAddr(), r); err != nil {
		panic(err)
	}
}
