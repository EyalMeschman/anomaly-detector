package validator

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
	tUsersInfoPath = "/users/info"
	tValidatePath  = "/validate"
)

var (
	tModel = &models.APIModel{
		Path:   tUsersInfoPath,
		Method: http.MethodGet,
		QueryParams: []*models.Parameter{
			{Name: "with_extra_data", Types: []models.ParamType{models.TypeBoolean}, Required: false},
		},
		Headers: []*models.Parameter{
			{Name: "Authorization", Types: []models.ParamType{models.TypeAuthToken, models.TypeUUID}, Required: true},
		},
	}
)

func TestValidateRequests(t *testing.T) {
	// Setup: Create store and add a test model
	tStoreMock := store.NewMockIModelStore(t)
	tValidatorMock := NewMockIRequestValidator(t)

	tHandler := &ValidateHandler{
		store:     tStoreMock,
		validator: tValidatorMock,
	}

	t.Run("success validating valid request", func(t *testing.T) {
		ctx := context.Background()

		request := models.Request{
			Path:   tUsersInfoPath,
			Method: http.MethodGet,
			QueryParams: []*models.RequestParam{
				{Name: "with_extra_data", Value: false},
			},
			Headers: []*models.RequestParam{
				{Name: "Authorization", Value: "Bearer abc123"},
			},
		}

		body, _ := json.Marshal(request)
		httpRequest := httptest.NewRequest(http.MethodPost, tValidatePath, bytes.NewReader(body))
		tRecorder := httptest.NewRecorder()

		tStoreMock.EXPECT().
			Get(ctx, tUsersInfoPath, http.MethodGet).
			Return(tModel, nil).Once()

		tValidatorMock.EXPECT().
			Validate(ctx, &request, tModel).
			Return(nil).Once()

		tHandler.Handle(tRecorder, httpRequest)

		assert.Equal(t, http.StatusOK, tRecorder.Code)

		var result models.ValidationResult

		expectedValidationResult := models.ValidationResult{
			Valid: true,
		}

		err := json.NewDecoder(tRecorder.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, result, expectedValidationResult)
	})

	t.Run("error when model not found", func(t *testing.T) {
		ctx := context.Background()

		request := models.Request{
			Path:   "/nonexistent",
			Method: http.MethodGet,
		}

		body, _ := json.Marshal(request)
		httpRequest := httptest.NewRequest(http.MethodPost, tValidatePath, bytes.NewReader(body))
		tRecorder := httptest.NewRecorder()

		tStoreMock.EXPECT().
			Get(ctx, "/nonexistent", http.MethodGet).
			Return(nil, assert.AnError).Once()

		tHandler.Handle(tRecorder, httpRequest)

		assert.Equal(t, http.StatusNotFound, tRecorder.Code)

		var response map[string]any

		err := json.NewDecoder(tRecorder.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "no model found for endpoint", response["error"])
	})

	t.Run("error with invalid JSON", func(t *testing.T) {
		// Create request with malformed JSON
		httpRequest := httptest.NewRequest(http.MethodPost, tValidatePath, bytes.NewReader([]byte("invalid json")))
		tRecorder := httptest.NewRecorder()

		tHandler.Handle(tRecorder, httpRequest)

		assert.Equal(t, http.StatusBadRequest, tRecorder.Code)

		var response map[string]any

		err := json.NewDecoder(tRecorder.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid JSON", response["error"])
	})
}
