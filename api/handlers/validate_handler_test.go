package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"anomaly_detector/models"
	"anomaly_detector/store"

	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestValidateRequests(t *testing.T) {
	// Setup: Create store and add a test model
	tStore := store.NewModelStore()
	ctx := context.Background()

	tModel := &models.APIModel{
		Path:   "/users/info",
		Method: "GET",
		QueryParams: []models.Parameter{
			{Name: "with_extra_data", Types: []string{"Boolean"}, Required: false},
		},
		Headers: []models.Parameter{
			{Name: "Authorization", Types: []string{"Auth-Token", "UUID"}, Required: true},
		},
		Body: []models.Parameter{},
	}
	_ = tStore.Store(ctx, tModel)

	t.Run("success validating valid request", func(t *testing.T) {
		tHandler := NewValidateHandler(tStore)

		// Create valid request
		tRequest := models.Request{
			Path:   "/users/info",
			Method: "GET",
			QueryParams: []models.RequestParam{
				{Name: "with_extra_data", Value: false},
			},
			Headers: []models.RequestParam{
				{Name: "Authorization", Value: "Bearer abc123"},
			},
			Body: []models.RequestParam{},
		}

		tBody, _ := json.Marshal(tRequest)
		tHTTPRequest := httptest.NewRequest("POST", "/validate", bytes.NewReader(tBody))
		tHTTPRequest.Header.Set("Content-Type", "application/json")

		tRecorder := httptest.NewRecorder()

		// Call handler
		tHandler.ValidateRequests(tRecorder, tHTTPRequest)

		// Assert response
		assert.Equal(t, http.StatusOK, tRecorder.Code)

		var tResult models.ValidationResult

		err := json.NewDecoder(tRecorder.Body).Decode(&tResult)
		assert.NoError(t, err)
		assert.False(t, tResult.IsAnomalous)
		assert.Empty(t, tResult.AnomalousFields)
	})

	t.Run("detects missing required parameter", func(t *testing.T) {
		tHandler := NewValidateHandler(tStore)

		// Create request missing required Authorization header
		tRequest := models.Request{
			Path:        "/users/info",
			Method:      "GET",
			QueryParams: []models.RequestParam{},
			Headers:     []models.RequestParam{}, // Missing required Authorization
			Body:        []models.RequestParam{},
		}

		tBody, _ := json.Marshal(tRequest)
		tHTTPRequest := httptest.NewRequest("POST", "/validate", bytes.NewReader(tBody))
		tRecorder := httptest.NewRecorder()

		// Call handler
		tHandler.ValidateRequests(tRecorder, tHTTPRequest)

		// Assert response
		assert.Equal(t, http.StatusOK, tRecorder.Code)

		var tResult models.ValidationResult

		err := json.NewDecoder(tRecorder.Body).Decode(&tResult)
		assert.NoError(t, err)
		assert.True(t, tResult.IsAnomalous)
		assert.Len(t, tResult.AnomalousFields, 1)
		assert.Contains(t, tResult.AnomalousFields[0].Reason, "missing")
	})

	t.Run("detects type mismatch", func(t *testing.T) {
		tHandler := NewValidateHandler(tStore)

		// Create request with wrong type for Authorization (should be Auth-Token or UUID)
		tRequest := models.Request{
			Path:   "/users/info",
			Method: "GET",
			Headers: []models.RequestParam{
				{Name: "Authorization", Value: float64(12345)}, // Wrong type
			},
			Body: []models.RequestParam{},
		}

		tBody, _ := json.Marshal(tRequest)
		tHTTPRequest := httptest.NewRequest("POST", "/validate", bytes.NewReader(tBody))
		tRecorder := httptest.NewRecorder()

		// Call handler
		tHandler.ValidateRequests(tRecorder, tHTTPRequest)

		// Assert response
		assert.Equal(t, http.StatusOK, tRecorder.Code)

		var tResult models.ValidationResult

		err := json.NewDecoder(tRecorder.Body).Decode(&tResult)
		assert.NoError(t, err)
		assert.True(t, tResult.IsAnomalous)
		assert.Contains(t, tResult.AnomalousFields[0].Reason, "type mismatch")
	})

	t.Run("error when model not found", func(t *testing.T) {
		tHandler := NewValidateHandler(tStore)

		// Create request for non-existent endpoint
		tRequest := models.Request{
			Path:   "/nonexistent",
			Method: "GET",
		}

		tBody, _ := json.Marshal(tRequest)
		tHTTPRequest := httptest.NewRequest("POST", "/validate", bytes.NewReader(tBody))
		tRecorder := httptest.NewRecorder()

		// Call handler
		tHandler.ValidateRequests(tRecorder, tHTTPRequest)

		// Assert response
		assert.Equal(t, http.StatusNotFound, tRecorder.Code)

		var tResponse map[string]any

		err := json.NewDecoder(tRecorder.Body).Decode(&tResponse)
		assert.NoError(t, err)
		assert.Equal(t, "error", tResponse["status"])
	})

	t.Run("error with invalid JSON", func(t *testing.T) {
		tHandler := NewValidateHandler(tStore)

		// Create request with malformed JSON
		tHTTPRequest := httptest.NewRequest("POST", "/validate", bytes.NewReader([]byte("invalid json")))
		tRecorder := httptest.NewRecorder()

		// Call handler
		tHandler.ValidateRequests(tRecorder, tHTTPRequest)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, tRecorder.Code)

		var tResponse map[string]any

		err := json.NewDecoder(tRecorder.Body).Decode(&tResponse)
		assert.NoError(t, err)
		assert.Equal(t, "error", tResponse["status"])
	})
}
