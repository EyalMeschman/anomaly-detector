package validator

import (
	"encoding/json"
	"fmt"
	"net/http"

	"anomaly_detector/api"
	"anomaly_detector/models"
	"anomaly_detector/store"
)

type IValidateHandler interface {
	api.IHandler
}

type validateHandler struct {
	store     store.IModelStore
	validator IRequestValidator
}

func NewValidateHandler(store store.IModelStore) IValidateHandler {
	return &validateHandler{
		store:     store,
		validator: NewRequestValidator(),
	}
}

func (h *validateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondError(w, http.StatusBadRequest, "invalid JSON provided")
		return
	}

	model, err := h.store.Get(ctx, req.Path, req.Method)
	if err != nil {
		api.RespondError(w, http.StatusNotFound, fmt.Sprintf("no model found for endpoint %s %s", req.Method, req.Path))
		return
	}

	result := h.validator.Validate(ctx, &req, model)

	validationResult := models.ValidationResult{
		Valid:     len(result) == 0,
		Anomalies: result,
	}

	api.RespondJSON(w, http.StatusOK, validationResult)
}
