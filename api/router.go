package api

import (
	"net/http"

	"anomaly_detector/api/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(
	modelsHandler *handlers.ModelsHandler,
	validateHandler *handlers.ValidateHandler,
) *mux.Router {
	r := mux.NewRouter()

	// Register routes
	r.HandleFunc("/models", modelsHandler.StoreAllModels).Methods("POST")
	r.HandleFunc("/validate", validateHandler.ValidateRequests).Methods("POST")

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	return r
}
