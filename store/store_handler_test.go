package store

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"anomaly_detector/models"

	"github.com/stretchr/testify/assert"
)

const (
	tModelsPath    = "/models"
	tUsersInfoPath = "/users/info"
)

var (
	tApiModels = []*models.APIModel{
		{
			Path:   tUsersInfoPath,
			Method: http.MethodGet,
			QueryParams: []*models.Parameter{
				{Name: "id", Types: []models.ParamType{models.TypeInt}, Required: true},
			},
		},
	}
)

func TestStoreAllModels(t *testing.T) {
	tStoreMock := NewMockIModelStore(t)

	tHandler := &storeHandler{
		store: tStoreMock,
	}

	t.Run("success storing models", func(t *testing.T) {
		ctx := context.Background()
		body, _ := json.Marshal(tApiModels)
		httpRequest := httptest.NewRequest(http.MethodPost, tModelsPath, bytes.NewReader(body))
		tRecorder := httptest.NewRecorder()

		tStoreMock.EXPECT().
			StoreAll(ctx, tApiModels).
			Return(true, nil).Once()

		tHandler.Handle(tRecorder, httpRequest)

		assert.Equal(t, http.StatusOK, tRecorder.Code)

		var response map[string]any

		err := json.NewDecoder(tRecorder.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "models stored successfully", response["message"])
	})

	t.Run("error with invalid JSON", func(t *testing.T) {
		httpRequest := httptest.NewRequest(http.MethodPost, tModelsPath, bytes.NewReader([]byte("invalid json")))
		tRecorder := httptest.NewRecorder()

		tHandler.Handle(tRecorder, httpRequest)

		assert.Equal(t, http.StatusBadRequest, tRecorder.Code)

		var response map[string]any

		err := json.NewDecoder(tRecorder.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid JSON", response["error"])
	})

	t.Run("bad request error", func(t *testing.T) {
		ctx := context.Background()
		body, _ := json.Marshal(tApiModels)
		httpRequest := httptest.NewRequest(http.MethodPost, tModelsPath, bytes.NewReader(body))
		tRecorder := httptest.NewRecorder()

		tStoreMock.EXPECT().
			StoreAll(ctx, tApiModels).
			Return(true, assert.AnError).Once()

		tHandler.Handle(tRecorder, httpRequest)

		assert.Equal(t, http.StatusBadRequest, tRecorder.Code)

		var response map[string]any

		err := json.NewDecoder(tRecorder.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "assert.AnError general error for testing", response["error"])
	})

	t.Run("internal server error", func(t *testing.T) {
		ctx := context.Background()
		body, _ := json.Marshal(tApiModels)
		httpRequest := httptest.NewRequest(http.MethodPost, tModelsPath, bytes.NewReader(body))
		tRecorder := httptest.NewRecorder()

		tStoreMock.EXPECT().
			StoreAll(ctx, tApiModels).
			Return(false, assert.AnError).Once()

		tHandler.Handle(tRecorder, httpRequest)

		assert.Equal(t, http.StatusInternalServerError, tRecorder.Code)

		var response map[string]any

		err := json.NewDecoder(tRecorder.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "internal server error", response["error"])
	})
}
