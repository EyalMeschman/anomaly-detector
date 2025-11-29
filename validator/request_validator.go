package validator

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"anomaly_detector/models"
)

const (
	cFieldQueryParams = "query_params"
	cFieldHeaders     = "headers"
	cFieldBody        = "body"
)

type IRequestValidator interface {
	Validate(ctx context.Context, req *models.Request, model *models.APIModel) []*models.FieldAnomaly
}

type requestValidator struct{}

func NewRequestValidator() IRequestValidator {
	return &requestValidator{}
}

func (rv *requestValidator) Validate(
	ctx context.Context, req *models.Request, model *models.APIModel) []*models.FieldAnomaly {
	slog.DebugContext(ctx, "Starting request validation",
		"path", req.Path,
		"method", req.Method,
	)

	var (
		anomalies []*models.FieldAnomaly
		wg        sync.WaitGroup
	)

	wg.Go(func() {
		anomalies = append(anomalies, rv.validateParameters(req.QueryParams, model.QueryParams, cFieldQueryParams)...)
	})
	wg.Go(func() {
		anomalies = append(anomalies, rv.validateParameters(req.Headers, model.Headers, cFieldHeaders)...)
	})
	wg.Go(func() {
		anomalies = append(anomalies, rv.validateParameters(req.Body, model.Body, cFieldBody)...)
	})

	wg.Wait()

	return anomalies
}

func (rv *requestValidator) validateParameters(
	requestParams []*models.RequestParam,
	modelParams []*models.Parameter,
	field string,
) []*models.FieldAnomaly {
	var anomalies []*models.FieldAnomaly

	// Build map of request parameters for quick lookup
	requestMap := make(map[string]any, len(requestParams))
	for _, rp := range requestParams {
		requestMap[rp.Name] = rp.Value
	}

	for _, modelParam := range modelParams {
		value, exists := requestMap[modelParam.Name]
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

		typeMatch := false

		for _, typeName := range modelParam.Types {
			if validateType(value, typeName) {
				typeMatch = true
				break
			}
		}

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
