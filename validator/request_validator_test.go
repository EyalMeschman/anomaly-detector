package validator

import (
	"context"
	"testing"

	"anomaly_detector/models"

	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestRequestValidator_Validate(t *testing.T) {
	t.Run("valid request with all required params", func(t *testing.T) {
		ctx := context.Background()
		validator := NewRequestValidator()

		tModel := &models.APIModel{
			Path:   "/test",
			Method: "GET",
			QueryParams: []models.Parameter{
				{Name: "id", Types: []string{"Int"}, Required: true},
			},
			Headers: []models.Parameter{
				{Name: "Auth", Types: []string{"String"}, Required: true},
			},
			Body: []models.Parameter{},
		}

		tRequest := &models.Request{
			Path:   "/test",
			Method: "GET",
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

	t.Run("missing required parameter", func(t *testing.T) {
		ctx := context.Background()
		validator := NewRequestValidator()

		tModel := &models.APIModel{
			Path:   "/test",
			Method: "GET",
			QueryParams: []models.Parameter{
				{Name: "id", Types: []string{"Int"}, Required: true},
			},
		}

		tRequest := &models.Request{
			Path:        "/test",
			Method:      "GET",
			QueryParams: []models.RequestParam{}, // Missing required "id"
		}

		result := validator.Validate(ctx, tRequest, tModel)
		assert.True(t, result.IsAnomalous)
		assert.Len(t, result.AnomalousFields, 1)
		assert.Contains(t, result.AnomalousFields[0].Reason, "missing")
		assert.Equal(t, "id", result.AnomalousFields[0].Name)
		assert.Equal(t, "query_params", result.AnomalousFields[0].Location)
	})

	t.Run("type mismatch", func(t *testing.T) {
		ctx := context.Background()
		validator := NewRequestValidator()

		tModel := &models.APIModel{
			Path:   "/test",
			Method: "POST",
			Body: []models.Parameter{
				{Name: "id", Types: []string{"Int"}, Required: true},
			},
		}

		tRequest := &models.Request{
			Path:   "/test",
			Method: "POST",
			Body: []models.RequestParam{
				{Name: "id", Value: "string value"}, // Expected Int, got String
			},
		}

		result := validator.Validate(ctx, tRequest, tModel)
		assert.True(t, result.IsAnomalous)
		assert.Len(t, result.AnomalousFields, 1)
		assert.Contains(t, result.AnomalousFields[0].Reason, "type mismatch")
		assert.Equal(t, "id", result.AnomalousFields[0].Name)
	})

	t.Run("multiple types allowed - matches one", func(t *testing.T) {
		ctx := context.Background()
		validator := NewRequestValidator()

		tModel := &models.APIModel{
			Path:   "/test",
			Method: "GET",
			QueryParams: []models.Parameter{
				{Name: "identifier", Types: []string{"Int", "UUID"}, Required: true},
			},
		}

		tRequest := &models.Request{
			Path:   "/test",
			Method: "GET",
			QueryParams: []models.RequestParam{
				{Name: "identifier", Value: "46da6390-7c78-4a1c-9efa-7c0396067ce4"}, // Valid UUID
			},
		}

		result := validator.Validate(ctx, tRequest, tModel)
		assert.False(t, result.IsAnomalous)
		assert.Len(t, result.AnomalousFields, 0)
	})

	t.Run("optional parameter missing - valid", func(t *testing.T) {
		ctx := context.Background()
		validator := NewRequestValidator()

		tModel := &models.APIModel{
			Path:   "/test",
			Method: "GET",
			QueryParams: []models.Parameter{
				{Name: "optional_param", Types: []string{"String"}, Required: false},
			},
		}

		tRequest := &models.Request{
			Path:        "/test",
			Method:      "GET",
			QueryParams: []models.RequestParam{}, // Missing optional param is OK
		}

		result := validator.Validate(ctx, tRequest, tModel)
		assert.False(t, result.IsAnomalous)
		assert.Len(t, result.AnomalousFields, 0)
	})

	t.Run("multiple anomalies", func(t *testing.T) {
		ctx := context.Background()
		validator := NewRequestValidator()

		tModel := &models.APIModel{
			Path:   "/test",
			Method: "POST",
			Headers: []models.Parameter{
				{Name: "Authorization", Types: []string{"Auth-Token"}, Required: true},
			},
			Body: []models.Parameter{
				{Name: "email", Types: []string{"Email"}, Required: true},
			},
		}

		tRequest := &models.Request{
			Path:    "/test",
			Method:  "POST",
			Headers: []models.RequestParam{}, // Missing Authorization
			Body:    []models.RequestParam{}, // Missing email
		}

		result := validator.Validate(ctx, tRequest, tModel)
		assert.True(t, result.IsAnomalous)
		assert.Len(t, result.AnomalousFields, 2)

		// Check both anomalies are present
		locations := make(map[string]bool)
		for _, anomaly := range result.AnomalousFields {
			locations[anomaly.Location] = true
			assert.Contains(t, anomaly.Reason, "missing")
		}

		assert.True(t, locations["headers"])
		assert.True(t, locations["body"])
	})
}
