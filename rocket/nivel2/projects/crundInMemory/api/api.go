package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"

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
			r.Post("/", handlerInsertUser(db))
			r.Get("/", handlerFindAllUser(db))
			r.Get("/{idUser}", handlerFindByIdUser(db))
			r.Put("/{idUser}", handlerUpdateUser(db))
			r.Delete("/{idUser}", handlerDeleteUser(db))
		})
	})

	return r
}

func handlerInsertUser(db map[string]Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body Users
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(w, apiResponse{Error: "invalid body"}, http.StatusUnprocessableEntity)
			return
		}

		if _, err := url.Parse(body.FirstName); err != nil {
			sendJSON(w, apiResponse{Error: "invalid first name passed"}, http.StatusBadRequest)
		}
		if _, err := url.Parse(body.LastName); err != nil {
			sendJSON(w, apiResponse{Error: "invalid last name passed"}, http.StatusBadRequest)
		}
		if _, err := url.Parse(body.Biography); err != nil {
			sendJSON(w, apiResponse{Error: "invalid biography passed"}, http.StatusBadRequest)
		}

		_uuid := uuid.New()

		user := Users{
			ID:        _uuid.String(),
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Biography: body.Biography,
		}

		db[_uuid.String()] = user

		sendJSON(w, apiResponse{Data: db[_uuid.String()]}, http.StatusCreated)
	}
}

func handlerFindAllUser(db map[string]Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := []Users{}

		for _, user := range db {
			users = append(users, user)
		}

		sendJSON(w, apiResponse{Data: users}, http.StatusOK)
	}
}

func handlerFindByIdUser(db map[string]Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "idUser")

		data, ok := db[id]

		if !ok {
			sendJSON(w, apiResponse{Error: "invalid id passed"}, http.StatusNotFound)
		}

		sendJSON(w, apiResponse{Data: data}, http.StatusOK)
	}
}

func handlerUpdateUser(db map[string]Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "idUser")

		data, ok := db[id]
		if !ok {
			sendJSON(w, apiResponse{Error: "invalid id passed"}, http.StatusNotFound)
		}

		var body Users
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(w, apiResponse{Error: "invalid body"}, http.StatusUnprocessableEntity)
			return
		}

		if _, err := url.Parse(body.FirstName); err != nil {
			sendJSON(w, apiResponse{Error: "invalid first name passed"}, http.StatusBadRequest)
		}
		if _, err := url.Parse(body.LastName); err != nil {
			sendJSON(w, apiResponse{Error: "invalid last name passed"}, http.StatusBadRequest)
		}
		if _, err := url.Parse(body.Biography); err != nil {
			sendJSON(w, apiResponse{Error: "invalid biography passed"}, http.StatusBadRequest)
		}

		data.FirstName = body.FirstName
		data.LastName = body.LastName
		data.Biography = body.Biography

		db[id] = data

		sendJSON(w, apiResponse{Data: db[id]}, http.StatusOK)
	}
}

func handlerDeleteUser(db map[string]Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "idUser")
		delete(db, id)

		sendJSON(w, apiResponse{Data: "successful"}, http.StatusOK)
	}
}
