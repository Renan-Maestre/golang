package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func HandlerFindAllUser(db map[string]users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := []Users{}

		for _, user := range db {
			users = append(users, user)
		}

		sendJSON(w, apiResponse{Data: users}, http.StatusOK)
	}
}

func HandlerFindByIdUser(db map[string]Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "idUser")

		data, ok := db[id]

		if !ok {
			sendJSON(w, apiResponse{Error: "invalid id passed"}, http.StatusNotFound)
		}

		sendJSON(w, apiResponse{Data: data}, http.StatusOK)
	}
}
