package validator

import (
	"context"
	"fmt"
	"log/slog"

	"anomaly_detector/models"
)

type IRequestValidator interface {
	Validate(ctx context.Context, req *models.Request, model *models.APIModel) []*models.FieldAnomaly
}

type RequestValidator struct{}

func NewRequestValidator() IRequestValidator {
	return &RequestValidator{}
}

func (v *RequestValidator) Validate(
	ctx context.Context, req *models.Request, model *models.APIModel) []*models.FieldAnomaly {
	slog.DebugContext(ctx, "Starting request validation",
		"path", req.Path,
		"method", req.Method,
	)

	anomalies := []*models.FieldAnomaly{}

	// Validate query parameters
	anomalies = append(anomalies, v.validateParameters(req.QueryParams, model.QueryParams, "query_params")...)

	// Validate headers
	anomalies = append(anomalies, v.validateParameters(req.Headers, model.Headers, "headers")...)

	// Validate body
	anomalies = append(anomalies, v.validateParameters(req.Body, model.Body, "body")...)

	return anomalies
}

func (v *RequestValidator) validateParameters(
	requestParams []*models.RequestParam,
	modelParams []*models.Parameter,
	field string,
) []*models.FieldAnomaly {
	anomalies := []*models.FieldAnomaly{}

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
				anomalies = append(anomalies, &models.FieldAnomaly{
					Field:         field,
					ParameterName: modelParam.Name,
					Reason:        fmt.Sprintf("required parameter %q is missing", modelParam.Name),
				})
			}

			continue
		}

		// Try to match against any of the allowed types
		typeMatch := false

		for _, typeName := range modelParam.Types {
			if validateType(value, typeName) {
				typeMatch = true
				break
			}
		}

		// If no type matched, add anomaly
		if !typeMatch {
			anomalies = append(anomalies, &models.FieldAnomaly{
				Field:         field,
				ParameterName: modelParam.Name,
				Reason:        fmt.Sprintf("type mismatch: expected one of %v types, but got the type %T", modelParam.Types, value),
			})
		}
	}

	return anomalies
}
