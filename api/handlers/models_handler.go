package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"anomaly_detector/models"
	"anomaly_detector/store"
)

// ModelsHandler handles API model storage requests
type ModelsHandler struct {
	store store.IModelStore
}

// NewModelsHandler creates a new ModelsHandler
func NewModelsHandler(store store.IModelStore) *ModelsHandler {
	return &ModelsHandler{store: store}
}

func (h *ModelsHandler) StoreAllModels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode request body
	var apiModels []*models.APIModel
	if err := json.NewDecoder(r.Body).Decode(&apiModels); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON", err)
		return
	}

	// Store all models
	storedCount, err := h.store.StoreAll(ctx, apiModels)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to store models", err)
		return
	}

	slog.InfoContext(ctx, "models stored successfully", "count", storedCount)

	// Respond with success
	response := map[string]any{
		"status":        "success",
		"models_stored": storedCount,
	}
	respondJSON(w, http.StatusOK, response)
}
