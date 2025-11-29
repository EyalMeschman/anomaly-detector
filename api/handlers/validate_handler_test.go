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
	"anomaly_detector/validator"

	"github.com/stretchr/testify/assert"
)

const (
	tUsersInfoPath = "/users/info"
	tValidatePath  = "/validate"
	tGetMethod     = "GET"
	tPostMethod    = "POST"
)

var (
	tModel = &models.APIModel{
		Path:   tUsersInfoPath,
		Method: tGetMethod,
		QueryParams: []models.Parameter{
			{Name: "with_extra_data", Types: []string{models.TypeBoolean}, Required: false},
		},
		Headers: []models.Parameter{
			{Name: "Authorization", Types: []string{models.TypeAuthToken, models.TypeUUID}, Required: true},
		},
		Body: []models.Parameter{},
	}
)

func TestValidateRequests(t *testing.T) {
	// Setup: Create store and add a test model
	tStoreMock := store.NewMockIModelStore(t)
	tValidatorMock := validator.NewMockIRequestValidator(t)

	tHandler := &ValidateHandler{
		store:     tStoreMock,
		validator: tValidatorMock,
	}

	t.Run("success validating valid request", func(t *testing.T) {
		ctx := context.Background()

		request := models.Request{
			Path:   tUsersInfoPath,
			Method: tGetMethod,
			QueryParams: []models.RequestParam{
				{Name: "with_extra_data", Value: false},
			},
			Headers: []models.RequestParam{
				{Name: "Authorization", Value: "Bearer abc123"},
			},
			Body: []models.RequestParam{},
		}

		body, _ := json.Marshal(request)
		httpRequest := httptest.NewRequest(tPostMethod, tValidatePath, bytes.NewReader(body))
		tRecorder := httptest.NewRecorder()

		tStoreMock.EXPECT().
			Get(ctx, tUsersInfoPath, tGetMethod).
			Return(tModel, nil).Once()

		tValidatorMock.EXPECT().
			Validate(ctx, &request, tModel).
			Return(&models.ValidationResult{
				IsAnomalous:     false,
				AnomalousFields: []models.FieldAnomaly{},
			}).Once()

		tHandler.ValidateRequests(tRecorder, httpRequest)

		assert.Equal(t, http.StatusOK, tRecorder.Code)

		var result models.ValidationResult

		_ = json.NewDecoder(tRecorder.Body).Decode(&result)
		assert.False(t, result.IsAnomalous)
		assert.Empty(t, result.AnomalousFields)
	})

	t.Run("error when model not found", func(t *testing.T) {
		ctx := context.Background()

		request := models.Request{
			Path:   "/nonexistent",
			Method: tGetMethod,
		}

		body, _ := json.Marshal(request)
		httpRequest := httptest.NewRequest(tPostMethod, tValidatePath, bytes.NewReader(body))
		tRecorder := httptest.NewRecorder()

		tStoreMock.EXPECT().
			Get(ctx, "/nonexistent", tGetMethod).
			Return(nil, assert.AnError).Once()

		tHandler.ValidateRequests(tRecorder, httpRequest)

		assert.Equal(t, http.StatusNotFound, tRecorder.Code)

		var response map[string]any

		_ = json.NewDecoder(tRecorder.Body).Decode(&response)
		assert.Equal(t, "error", response["status"])
	})

	t.Run("error with invalid JSON", func(t *testing.T) {
		// Create request with malformed JSON
		httpRequest := httptest.NewRequest(tPostMethod, tValidatePath, bytes.NewReader([]byte("invalid json")))
		tRecorder := httptest.NewRecorder()

		tHandler.ValidateRequests(tRecorder, httpRequest)

		assert.Equal(t, http.StatusBadRequest, tRecorder.Code)

		var response map[string]any

		_ = json.NewDecoder(tRecorder.Body).Decode(&response)
		assert.Equal(t, "error", response["status"])
	})
}
