package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"

	"projetoComDb/internal/store"
)

type shortenURLRequest struct {
	URL string `json:"url"`
}

type shortenURLResponse struct {
	Code string `json:"code"`
}

func handleShortenURL(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body shortenURLRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			SendJSON(w, apiResponse{Error: "invalid body"}, http.StatusUnprocessableEntity)
			return
		}

		if _, err := url.Parse(body.URL); err != nil {
			SendJSON(w, apiResponse{Error: "invalid url passed"}, http.StatusBadRequest)
		}

		code, err := store.SaveShortenedURL(r.Context(), body.URL)
		if err != nil {
			slog.Error("failed to create code ", "error", err)
			SendJSON(w, apiResponse{Error: "something went wrong"}, http.StatusInternalServerError)
			return
		}

		SendJSON(w, apiResponse{Data: shortenURLResponse{Code: code}}, http.StatusCreated)
	}
}
