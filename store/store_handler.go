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

type storeHandler struct {
	store IModelStore
}

func NewStoreHandler(store IModelStore) IStoreHandler {
	return &storeHandler{store: store}
}

func (h *storeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var apiModels []*models.APIModel
	if err := json.NewDecoder(r.Body).Decode(&apiModels); err != nil {
		api.RespondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	ok, err := h.store.StoreAll(ctx, apiModels)
	if err != nil {
		if !ok {
			// Internal errors - log but don't expose details to prevent information leakage
			slog.ErrorContext(ctx, "error storing models", "error", err)
			api.RespondError(w, http.StatusInternalServerError, "internal server error")

			return
		}

		// Safe to expose validation errors to users
		api.RespondError(w, http.StatusBadRequest, err.Error())

		return
	}

	response := map[string]any{
		"message": "models stored successfully",
	}
	api.RespondJSON(w, http.StatusOK, response)
}
