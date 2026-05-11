package routes

import (
	"github.com/FoPQer/go-shortener/internal/handlers"
	"github.com/FoPQer/go-shortener/internal/middlewares"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/go-chi/chi/v5"
)

// InitWebRoutes registers all HTTP routes and attaches shared middlewares.
func InitWebRoutes(r *chi.Mux, handler *handlers.Handler, dbHandler *handlers.DBHandler, authMiddleware *middlewares.AuthMiddleware) {
	base := service.GetBasePrefix()

	r.Use(middlewares.WithGzip, middlewares.WithLogging)

	auth := r.Group(func(auth chi.Router) {
		auth.Use(authMiddleware.WithAuth)
	})

	r.Get("/ping", dbHandler.GetPing)
	auth.Get(base+"{id}", handler.GetURL)
	auth.Post("/", handler.PostURL)

	auth.Route("/api", func(api chi.Router) {
		api.Post("/shorten", handler.PostURLByJSON)
		api.Post("/shorten/batch", handler.PostBatchURLByJSON)
		api.Route("/user", func(user chi.Router) {
			user.Get("/urls", handler.GetUserURLs)
			user.Delete("/urls", handler.DeleteUserURLs)
		})
	})
}
