package store

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"anomaly_detector/models"
)

type IModelStore interface {
	StoreAll(ctx context.Context, models []*models.APIModel) (bool, error)
	Get(ctx context.Context, path, method string) (*models.APIModel, error)
}

type modelStore struct {
	mu     sync.RWMutex
	models map[string]*models.APIModel
}

func NewModelStore() IModelStore {
	return &modelStore{
		models: make(map[string]*models.APIModel),
	}
}

func (s *modelStore) store(ctx context.Context, apiModel *models.APIModel) {
	key := getKey(apiModel.Path, apiModel.Method)

	s.models[key] = apiModel

	slog.InfoContext(ctx, "Model stored", "path", apiModel.Path, "method", apiModel.Method)
}

// StoreAll stores multiple API models and returns a bool indicating whether any error
// was caused by user input (true) or an internal server error (false).
// Currently, only user input errors are possible, but this may change in the future to support database storage.
func (s *modelStore) StoreAll(ctx context.Context, apiModels []*models.APIModel) (bool, error) {
	for _, model := range apiModels {
		if model == nil || model.Path == "" || model.Method == "" {
			return true, fmt.Errorf("invalid model in batch")
		}

		key := getKey(model.Path, model.Method)
		if _, exists := s.models[key]; exists {
			return true, fmt.Errorf("model already exists for path %s and method %s", model.Path, model.Method)
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, apiModel := range apiModels {
		s.store(ctx, apiModel)
	}

	return true, nil
}

func (s *modelStore) Get(ctx context.Context, path, method string) (*models.APIModel, error) {
	slog.InfoContext(ctx, "Getting model", "path", path, "method", method)

	s.mu.RLock()
	defer s.mu.RUnlock()

	key := getKey(path, method)

	model, exists := s.models[key]
	if !exists {
		return nil, fmt.Errorf("model not found for path %s and method %s", path, method)
	}

	slog.InfoContext(ctx, "Model retrieved", "path", path, "method", method)

	return model, nil
}

// getKey returns a unique identifier for this API model
func getKey(path, method string) string {
	return fmt.Sprintf("%s:%s", path, method)
}
