package validator

import (
	"context"
	"net/http"
	"testing"

	"anomaly_detector/models"

	"github.com/stretchr/testify/assert"
)

const (
	tTestPath = "/test"
)

func TestRequestValidator_Validate(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		ctx := context.Background()
		validator := NewRequestValidator()

		tModel := &models.APIModel{
			Path:   tTestPath,
			Method: http.MethodGet,
			QueryParams: []*models.Parameter{
				{Name: "id", Types: []models.ParamType{models.TypeInt}, Required: true},
			},
			Headers: []*models.Parameter{
				{Name: "Auth", Types: []models.ParamType{models.TypeString}, Required: true},
			},
		}

		tRequest := &models.Request{
			Path:   tTestPath,
			Method: http.MethodGet,
			QueryParams: []*models.RequestParam{
				{Name: "id", Value: float64(123)},
			},
			Headers: []*models.RequestParam{
				{Name: "Auth", Value: "Bearer token123"},
			},
		}

		result := validator.Validate(ctx, tRequest, tModel)
		assert.Empty(t, result)
	})

	t.Run("multiple anomalies", func(t *testing.T) {
		ctx := context.Background()
		validator := NewRequestValidator()

		tModel := &models.APIModel{
			Path:   tTestPath,
			Method: http.MethodPost,
			Headers: []*models.Parameter{
				{Name: "Authorization", Types: []models.ParamType{models.TypeAuthToken}, Required: true},
			},
			Body: []*models.Parameter{
				{Name: "id", Types: []models.ParamType{models.TypeInt}, Required: true},
			},
		}

		tRequest := &models.Request{
			Path:    tTestPath,
			Method:  http.MethodPost,
			Headers: []*models.RequestParam{}, // Missing Authorization
			Body: []*models.RequestParam{
				{Name: "id", Value: "string value"}, // Expected Int, got String
			},
		}

		expectedAnomalousFields := []*models.FieldAnomaly{
			{
				Field:         "headers",
				ParameterName: "Authorization",
				Reason:        "required parameter \"Authorization\" is missing",
			},
			{
				Field:         "body",
				ParameterName: "id",
				Reason:        "type mismatch: expected one of [Int] types, but got the type string",
			},
		}

		result := validator.Validate(ctx, tRequest, tModel)
		assert.Equal(t, expectedAnomalousFields, result)
	})
}
