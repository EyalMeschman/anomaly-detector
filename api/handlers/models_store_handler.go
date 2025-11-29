package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"anomaly_detector/models"
	"anomaly_detector/store"
)

type ModelsHandler struct {
	store store.IModelStore
}

func NewModelsHandler(store store.IModelStore) *ModelsHandler {
	return &ModelsHandler{store: store}
}

func (h *ModelsHandler) StoreAllModels(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var apiModels []*models.APIModel
	if err := json.NewDecoder(r.Body).Decode(&apiModels); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON", err)
		return
	}

	storedCount, err := h.store.StoreAll(ctx, apiModels)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to store models", err)
		return
	}

	slog.InfoContext(ctx, "models stored successfully", "count", storedCount)

	response := map[string]any{
		"status":        "success",
		"models_stored": storedCount,
	}
	respondJSON(w, http.StatusOK, response)
}
