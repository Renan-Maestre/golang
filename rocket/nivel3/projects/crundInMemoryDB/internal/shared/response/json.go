package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ApiResponse struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, resp ApiResponse, status int) {
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("failed to marshal json data", "error", err, "response", resp)
		http.Error(
			w,
			"something went wrong",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write response to client ", "error", err)

		return
	}
}
