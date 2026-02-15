package routes

import (
	"github.com/FoPQer/go-shortener/internal/handler"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

func InitWebRoutes(r *chi.Mux) {
	base := service.GetBasePrefix()

	r.Get(base+"{id}", handler.GetURL)
	r.Post("/", handler.PostURL)
}
