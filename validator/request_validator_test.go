package validator

import (
	"context"
	"testing"

	"anomaly_detector/models"

	"github.com/stretchr/testify/assert"
)

const (
	tTestPath   = "/test"
	tGetMethod  = "GET"
	tPostMethod = "POST"
)

func TestRequestValidator_Validate(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		ctx := context.Background()
		validator := NewRequestValidator()

		tModel := &models.APIModel{
			Path:   tTestPath,
			Method: tGetMethod,
			QueryParams: []models.Parameter{
				{Name: "id", Types: []string{models.TypeInt}, Required: true},
			},
			Headers: []models.Parameter{
				{Name: "Auth", Types: []string{models.TypeString}, Required: true},
			},
			Body: []models.Parameter{},
		}

		tRequest := &models.Request{
			Path:   tTestPath,
			Method: tGetMethod,
			QueryParams: []models.RequestParam{
				{Name: "id", Value: float64(123)},
			},
			Headers: []models.RequestParam{
				{Name: "Auth", Value: "Bearer token123"},
			},
			Body: []models.RequestParam{},
		}

		result := validator.Validate(ctx, tRequest, tModel)
		assert.False(t, result.IsAnomalous)
		assert.Len(t, result.AnomalousFields, 0)
	})

	t.Run("multiple anomalies", func(t *testing.T) {
		ctx := context.Background()
		validator := NewRequestValidator()

		tModel := &models.APIModel{
			Path:   tTestPath,
			Method: tPostMethod,
			Headers: []models.Parameter{
				{Name: "Authorization", Types: []string{models.TypeAuthToken}, Required: true},
			},
			Body: []models.Parameter{
				{Name: "id", Types: []string{models.TypeInt}, Required: true},
			},
		}

		tRequest := &models.Request{
			Path:    tTestPath,
			Method:  tPostMethod,
			Headers: []models.RequestParam{}, // Missing Authorization
			Body: []models.RequestParam{
				{Name: "id", Value: "string value"}, // Expected Int, got String
			},
		}

		expectedValidationResult := &models.ValidationResult{
			IsAnomalous: true,
			AnomalousFields: []models.FieldAnomaly{
				{
					Location: "headers",
					Name:     "Authorization",
					Reason:   "required parameter 'Authorization' is missing",
				},
				{
					Location: "body",
					Name:     "id",
					Reason:   "type mismatch: expected one of [Int], got string",
				},
			},
		}

		result := validator.Validate(ctx, tRequest, tModel)
		assert.Equal(t, expectedValidationResult, result)
	})
}
