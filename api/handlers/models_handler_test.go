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

const (
	tModelsPath = "/models"
)

var (
	tApiModels = []*models.APIModel{
		{
			Path:   tUsersInfoPath,
			Method: tGetMethod,
			QueryParams: []models.Parameter{
				{Name: "id", Types: []string{models.TypeInt}, Required: true},
			},
			Headers: []models.Parameter{},
			Body:    []models.Parameter{},
		},
	}
)

func TestStoreAllModels(t *testing.T) {
	tStoreMock := store.NewMockIModelStore(t)

	tHandler := &ModelsHandler{
		store: tStoreMock,
	}

	t.Run("success storing models", func(t *testing.T) {
		ctx := context.Background()
		body, _ := json.Marshal(tApiModels)
		httpRequest := httptest.NewRequest(tPostMethod, tModelsPath, bytes.NewReader(body))
		tRecorder := httptest.NewRecorder()

		tStoreMock.EXPECT().
			StoreAll(ctx, tApiModels).
			Return(2, nil).Once()

		tHandler.StoreAllModels(tRecorder, httpRequest)

		assert.Equal(t, http.StatusOK, tRecorder.Code)

		var response map[string]any

		_ = json.NewDecoder(tRecorder.Body).Decode(&response)
		assert.Equal(t, float64(2), response["models_stored"])
	})

	t.Run("error with invalid JSON", func(t *testing.T) {
		httpRequest := httptest.NewRequest(tPostMethod, tModelsPath, bytes.NewReader([]byte("invalid json")))
		tRecorder := httptest.NewRecorder()

		tHandler.StoreAllModels(tRecorder, httpRequest)

		assert.Equal(t, http.StatusBadRequest, tRecorder.Code)

		var response map[string]any

		_ = json.NewDecoder(tRecorder.Body).Decode(&response)
		assert.Equal(t, "error", response["status"])
	})

	t.Run("error when store fails", func(t *testing.T) {
		ctx := context.Background()
		body, _ := json.Marshal(tApiModels)
		httpRequest := httptest.NewRequest(tPostMethod, tModelsPath, bytes.NewReader(body))
		tRecorder := httptest.NewRecorder()

		tStoreMock.EXPECT().
			StoreAll(ctx, tApiModels).
			Return(0, assert.AnError).Once()

		tHandler.StoreAllModels(tRecorder, httpRequest)

		assert.Equal(t, http.StatusInternalServerError, tRecorder.Code)

		var response map[string]any

		_ = json.NewDecoder(tRecorder.Body).Decode(&response)
		assert.Equal(t, "error", response["status"])
	})
}
