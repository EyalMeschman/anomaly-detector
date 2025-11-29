package store

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"anomaly_detector/api"
	"anomaly_detector/models"
)

type IStoreHandler interface {
	api.IHandler
}

type StoreHandler struct {
	store IModelStore
}

func NewStoreHandler(store IModelStore) IStoreHandler {
	return &StoreHandler{store: store}
}

func (h *StoreHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var apiModels []*models.APIModel
	if err := json.NewDecoder(r.Body).Decode(&apiModels); err != nil {
		api.RespondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	err := h.store.StoreAll(ctx, apiModels)
	if err != nil {
		if isValidationError(err) {
			// Safe to expose validation errors to users
			api.RespondError(w, http.StatusBadRequest, err.Error())
		} else {
			// Internal errors - log but don't expose details to prevent information leakage
			slog.ErrorContext(ctx, "Internal error storing models", "error", err)
			api.RespondError(w, http.StatusInternalServerError, "internal server error")
		}

		return
	}

	response := map[string]any{
		"message": "models stored successfully",
	}
	api.RespondJSON(w, http.StatusOK, response)
}
