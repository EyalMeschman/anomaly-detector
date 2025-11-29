package store

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"anomaly_detector/models"
)

type IModelStore interface {
	StoreAll(ctx context.Context, models []*models.APIModel) error
	Get(ctx context.Context, path, method string) (*models.APIModel, error)
}

type ModelStore struct {
	mu     sync.RWMutex
	models map[string]*models.APIModel
}

func NewModelStore() IModelStore {
	return &ModelStore{
		models: make(map[string]*models.APIModel),
	}
}

func (s *ModelStore) store(ctx context.Context, apiModel *models.APIModel) {
	key := getApiModelKey(apiModel.Path, apiModel.Method)

	s.models[key] = apiModel

	slog.InfoContext(ctx, "Model stored", "path", apiModel.Path, "method", apiModel.Method)
}

func (s *ModelStore) StoreAll(ctx context.Context, apiModels []*models.APIModel) error {
	// Validate all first (without locking the whole operation)
	for _, model := range apiModels {
		if model == nil || model.Path == "" || model.Method == "" {
			return &ModelStoreError{ErrType: ErrorInvalidModel}
		}

		key := getApiModelKey(model.Path, model.Method)
		if _, exists := s.models[key]; exists {
			return &ModelStoreError{
				ErrType:        ErrorDuplicateModel,
				AdditionalInfo: []string{model.Method, model.Path},
			}
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, apiModel := range apiModels {
		s.store(ctx, apiModel)
	}

	return nil
}

func (s *ModelStore) Get(ctx context.Context, path, method string) (*models.APIModel, error) {
	slog.InfoContext(ctx, "Getting model", "path", path, "method", method)

	s.mu.RLock()
	defer s.mu.RUnlock()

	key := getApiModelKey(path, method)

	model, exists := s.models[key]
	if !exists {
		return nil, &ModelStoreError{
			ErrType:        ErrorModelNotFound,
			AdditionalInfo: []string{method, path},
		}
	}

	slog.InfoContext(ctx, "Model retrieved", "path", path, "method", method)

	return model, nil
}

// getApiModelKey returns a unique identifier for this API model
func getApiModelKey(path, method string) string {
	return fmt.Sprintf("%s:%s", path, method)
}
