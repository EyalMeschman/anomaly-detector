package validator

import (
	"context"
	"log/slog"

	"anomaly_detector/internal/models"
)

// IRequestValidator defines the interface for validating requests against API models
type IRequestValidator interface {
	Validate(ctx context.Context, req *models.Request, model *models.APIModel) *models.ValidationResult
}

// RequestValidator orchestrates validation of entire requests
type RequestValidator struct{}

// NewRequestValidator creates a new request validator
func NewRequestValidator() IRequestValidator {
	return &RequestValidator{}
}

// Validate validates a request against an API model
func (v *RequestValidator) Validate(
	ctx context.Context, req *models.Request, model *models.APIModel) *models.ValidationResult {
	slog.DebugContext(ctx, "Starting request validation",
		"path", req.Path,
		"method", req.Method,
	)
	result := &models.ValidationResult{
		IsAnomalous:     false,
		AnomalousFields: []models.FieldAnomaly{},
	}

	// Validate query parameters
	anomalies := v.validateParameters(req.QueryParams, model.QueryParams, "query_params")
	result.AnomalousFields = append(result.AnomalousFields, anomalies...)

	// Validate headers
	anomalies = v.validateParameters(req.Headers, model.Headers, "headers")
	result.AnomalousFields = append(result.AnomalousFields, anomalies...)

	// Validate body
	anomalies = v.validateParameters(req.Body, model.Body, "body")
	result.AnomalousFields = append(result.AnomalousFields, anomalies...)

	// Set anomalous flag if any anomalies found
	if len(result.AnomalousFields) > 0 {
		result.IsAnomalous = true
	}

	return result
}

// validateParameters validates a set of request parameters against model parameters
func (v *RequestValidator) validateParameters(
	requestParams []models.RequestParam,
	modelParams []models.Parameter,
	location string,
) []models.FieldAnomaly {
	anomalies := []models.FieldAnomaly{}

	// Build map of request parameters for quick lookup
	requestMap := make(map[string]any, len(requestParams))
	for _, rp := range requestParams {
		requestMap[rp.Name] = rp.Value
	}

	// Validate each model parameter
	for _, modelParam := range modelParams {
		value, exists := requestMap[modelParam.Name]

		// Check if required parameter is missing
		if !exists {
			if modelParam.Required {
				anomalies = append(anomalies, models.FieldAnomaly{
					Location: location,
					Name:     modelParam.Name,
					Reason:   "required parameter '" + modelParam.Name + "' is missing",
				})
			}

			continue
		}

		// Try to match against any of the allowed types
		typeMatch := false

		var lastError error

		for _, typeName := range modelParam.Types {
			err := validateType(value, typeName)
			if err == nil {
				typeMatch = true
				break
			}

			lastError = err
		}

		// If no type matched, add anomaly
		if !typeMatch {
			anomalies = append(anomalies, models.FieldAnomaly{
				Location: location,
				Name:     modelParam.Name,
				Reason:   "type mismatch: " + lastError.Error(),
			})
		}
	}

	return anomalies
}
