package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}

func respondError(w http.ResponseWriter, status int, message string, err error) {
	response := map[string]any{
		"status":  "error",
		"message": message,
	}
	if err != nil {
		response["error"] = err.Error()
	}

	respondJSON(w, status, response)
}
