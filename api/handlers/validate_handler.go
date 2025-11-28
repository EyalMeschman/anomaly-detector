package handlers

import (
	"encoding/json"
	"net/http"

	"anomaly_detector/models"
	"anomaly_detector/store"
	"anomaly_detector/validator"
)

type ValidateHandler struct {
	store     store.IModelStore
	validator validator.IRequestValidator
}

func NewValidateHandler(store store.IModelStore) *ValidateHandler {
	return &ValidateHandler{
		store:     store,
		validator: validator.NewRequestValidator(),
	}
}

func (h *ValidateHandler) ValidateRequests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON", err)
		return
	}

	model, err := h.store.Get(ctx, req.Path, req.Method)
	if err != nil {
		respondError(w, http.StatusNotFound, "no model found for endpoint", err)
		return
	}

	result := h.validator.Validate(ctx, &req, model)

	respondJSON(w, http.StatusOK, result)
}
