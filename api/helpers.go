package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}

func RespondError(w http.ResponseWriter, status int, message string) {
	response := map[string]any{
		"error": message,
	}

	RespondJSON(w, status, response)
}
