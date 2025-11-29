package validator

import (
	"encoding/json"
	"net/http"

	"anomaly_detector/api"
	"anomaly_detector/models"
	"anomaly_detector/store"
)

type IValidateHandler interface {
	api.IHandler
}

type ValidateHandler struct {
	store     store.IModelStore
	validator IRequestValidator
}

func NewValidateHandler(store store.IModelStore) IValidateHandler {
	return &ValidateHandler{
		store:     store,
		validator: NewRequestValidator(),
	}
}

func (h *ValidateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	model, err := h.store.Get(ctx, req.Path, req.Method)
	if err != nil {
		api.RespondError(w, http.StatusNotFound, "no model found for endpoint")
		return
	}

	result := h.validator.Validate(ctx, &req, model)

	validationResult := models.ValidationResult{
		Valid:     len(result) == 0,
		Anomalies: result,
	}

	api.RespondJSON(w, http.StatusOK, validationResult)
}
