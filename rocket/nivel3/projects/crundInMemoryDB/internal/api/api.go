package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type apiResponse struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

type id uuid.UUID

type Users struct {
	ID        string
	FirstName string
	LastName  string
	Biography string
}

func sendJSON(w http.ResponseWriter, resp apiResponse, status int) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("failed to marshal json data", "error", err)
		sendJSON(w, apiResponse{Error: "something went wrong"}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write responde to client", "error", err)
	}
}

func Newhandler(db map[string]Users) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", HandlerInsertUser(db))
			r.Get("/", HandlerFindAllUser(db))
			r.Get("/{idUser}", HandlerFindByIdUser(db))
			r.Put("/{idUser}", handlerUpdateUser(db))
			r.Delete("/{idUser}", handlerDeleteUser(db))
		})
	})

	return r
}

func handlerDeleteUser(db map[string]Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "idUser")
		delete(db, id)

		sendJSON(w, apiResponse{Data: "successful"}, http.StatusOK)
	}
}
