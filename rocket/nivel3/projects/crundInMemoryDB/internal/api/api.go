package api

import (
	"net/http"

	"crudMemoryDB/internal/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type id uuid.UUID

func NewRouter(userHandler *user.Handler) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.Create)
			r.Get("/", userHandler.FindAll)
			r.Get("/{idUser}", userHandler.FindByID)
			r.Put("/{idUser}", userHandler.Update)
			r.Delete("/{idUser}", userHandler.Delete)
		})
	})

	return r
}
