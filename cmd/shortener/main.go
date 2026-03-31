package main

import (
	"errors"
	"net/http"

	"github.com/FoPQer/go-shortener/internal/auth"
	"github.com/FoPQer/go-shortener/internal/config/db"
	"github.com/FoPQer/go-shortener/internal/config/flags"
	"github.com/FoPQer/go-shortener/internal/handlers"
	"github.com/FoPQer/go-shortener/internal/logger"
	"github.com/FoPQer/go-shortener/internal/middlewares"
	repoFactory "github.com/FoPQer/go-shortener/internal/repository/factory"
	"github.com/FoPQer/go-shortener/internal/routes"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	flags.ParseFlags()
	logger.InitLogger()
	pgxConf, err := db.InitPgsql()
	if errors.Is(err, db.ErrConnNotFound) {
		logger.GetSugar().Infoln("Database connection string not found, using file or memory repository")
	} else if err != nil {
		logger.GetSugar().Errorf("Error initializing database: %s", err.Error())
		panic(err)
	}
	if pgxConf.GetDBConn() != nil {
		defer pgxConf.GetDBConn().Close()
	} 

	factory := repoFactory.NewRepositoryFactory(pgxConf.GetDBConn(), service.GetFileStoragePath())
		
    urlRepo, err := factory.CreateUrlsRepository()
    if err != nil {
        panic(err)
    }

	userRepo, err := factory.CreateUserRepository()
    if err != nil {
        panic(err)
    }

    urlService := service.NewURLService(urlRepo)
	jsonService := service.NewJSONService()
	userService := service.NewUserService(userRepo)
	claimsService := auth.NewClaimsService()

	authMiddleware := middlewares.NewAuthMiddleware(userService, claimsService)

	handler := handlers.NewHandler(urlService, jsonService, userService)
	dbHandler := handlers.NewDBHandler(pgxConf.GetDBConn())


	r := chi.NewRouter()
	routes.InitWebRoutes(r, handler, dbHandler, authMiddleware)

	if err := http.ListenAndServe(service.GetRunAddr(), r); err != nil {
		panic(err)
	}
}
