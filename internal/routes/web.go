package routes

import (
	"github.com/FoPQer/go-shortener/internal/handlers"
	"github.com/FoPQer/go-shortener/internal/middlewares"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func InitWebRoutes(r *chi.Mux, urlService *service.URLService, jsonService *service.JSONService) {
	handler := handlers.NewHandler(urlService, jsonService)
	base := service.GetBasePrefix()

	r.Use(middlewares.WithGzip, middlewares.WithLogging)

	r.Get("/ping", handler.GetPing)
	r.Get(base+"{id}", handler.GetURL)
	r.Post("/", handler.PostURL)

	r.Route("/api", func(api chi.Router) {
		api.Post("/shorten", handler.PostURLByJSON)
		api.Post("/shorten/batch", handler.PostBatchURLByJSON)
	})
}
