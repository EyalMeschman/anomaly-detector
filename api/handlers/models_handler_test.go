package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"anomaly_detector/models"
	"anomaly_detector/store"

	"github.com/stretchr/testify/assert"
)

func TestStoreAllModels(t *testing.T) {
	t.Run("success storing valid models", func(t *testing.T) {
		// Create store and handler
		tStore := store.NewModelStore()
		tHandler := NewModelsHandler(tStore)

		// Create test models
		tModels := []*models.APIModel{
			{
				Path:   "/test1",
				Method: "GET",
				QueryParams: []models.Parameter{
					{Name: "param1", Types: []string{"String"}, Required: true},
				},
			},
			{
				Path:   "/test2",
				Method: "POST",
				Body: []models.Parameter{
					{Name: "field1", Types: []string{"Int"}, Required: true},
				},
			},
		}

		// Marshal to JSON
		tBody, err := json.Marshal(tModels)
		assert.NoError(t, err)

		// Create request
		tRequest := httptest.NewRequest("POST", "/models", bytes.NewReader(tBody))
		tRequest.Header.Set("Content-Type", "application/json")

		tRecorder := httptest.NewRecorder()

		// Call handler
		tHandler.StoreAllModels(tRecorder, tRequest)

		// Assert response
		assert.Equal(t, http.StatusOK, tRecorder.Code)

		var tResponse map[string]any

		err = json.NewDecoder(tRecorder.Body).Decode(&tResponse)
		assert.NoError(t, err)
		assert.Equal(t, "success", tResponse["status"])
		assert.Equal(t, float64(2), tResponse["models_stored"])
	})

	t.Run("error with invalid JSON", func(t *testing.T) {
		// Create handler
		tStore := store.NewModelStore()
		tHandler := NewModelsHandler(tStore)

		// Create request with invalid JSON
		tRequest := httptest.NewRequest("POST", "/models", bytes.NewReader([]byte("invalid json")))
		tRequest.Header.Set("Content-Type", "application/json")

		tRecorder := httptest.NewRecorder()

		// Call handler
		tHandler.StoreAllModels(tRecorder, tRequest)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, tRecorder.Code)

		var tResponse map[string]any

		err := json.NewDecoder(tRecorder.Body).Decode(&tResponse)
		assert.NoError(t, err)
		assert.Equal(t, "error", tResponse["status"])
		assert.Contains(t, tResponse["message"], "invalid JSON")
	})
}
