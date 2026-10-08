package user

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

func HandlerInsertUser(db map[string]Users) http.HandlerFunc {
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
