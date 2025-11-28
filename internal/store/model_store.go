package store

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"anomaly_detector/internal/models"
)

type IModelStore interface {
	Store(ctx context.Context, model *models.APIModel) error
	Get(ctx context.Context, path, method string) (*models.APIModel, error)
	GetAll(ctx context.Context) []*models.APIModel
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

func (s *ModelStore) Store(ctx context.Context, model *models.APIModel) error {
	if model == nil {
		return fmt.Errorf("model cannot be nil")
	}

	if model.Path == "" || model.Method == "" {
		return fmt.Errorf("path and method are required")
	}

	slog.InfoContext(ctx, "Storing model", "path", model.Path, "method", model.Method)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.models[model.Key()] = model

	slog.InfoContext(ctx, "Model stored", "path", model.Path, "method", model.Method)

	return nil
}

func (s *ModelStore) Get(ctx context.Context, path, method string) (*models.APIModel, error) {
	slog.InfoContext(ctx, "Getting model", "path", path, "method", method)

	s.mu.RLock()
	defer s.mu.RUnlock()

	key := path + ":" + method

	model, exists := s.models[key]
	if !exists {
		return nil, fmt.Errorf("model not found for %s %s", method, path)
	}

	slog.InfoContext(ctx, "Model retrieved", "path", path, "method", method)

	return model, nil
}

func (s *ModelStore) GetAll(ctx context.Context) []*models.APIModel {
	slog.InfoContext(ctx, "Getting all models")

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*models.APIModel, 0, len(s.models))
	for _, model := range s.models {
		result = append(result, model)
	}

	slog.InfoContext(ctx, "All models retrieved")

	return result
}
