package api

import (
	"errors"
	"log/slog"
	"net/http"

	"projetoComDb/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type getShortenedURLResponse struct {
	FullURL string `json:"full_url"`
}

func handleGetshortenedUrl(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		fullURL, err := store.GetFullURL(r.Context(), code)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				SendJSON(w, apiResponse{Error: "code not found"}, http.StatusNotFound)
				return
			}
			slog.Error("failed to get code ", "error", err)
			SendJSON(w, apiResponse{Error: "something went wrong"}, http.StatusInternalServerError)
			return
		}

		SendJSON(w, apiResponse{Data: getShortenedURLResponse{FullURL: fullURL}}, http.StatusOK)
	}
}
