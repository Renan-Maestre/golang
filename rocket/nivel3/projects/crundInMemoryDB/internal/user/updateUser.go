package user

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

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
